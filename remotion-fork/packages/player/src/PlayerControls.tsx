import type {MouseEventHandler, ReactNode, SyntheticEvent} from 'react';
import React, {useCallback, useEffect, useMemo, useRef, useState} from 'react';
import {Internals} from 'remotion';
import {DefaultPlayPauseButton} from './DefaultPlayPauseButton.js';
import {MediaVolumeSlider} from './MediaVolumeSlider.js';
import {PlaybackrateControl, playerButtonStyle} from './PlaybackrateControl.js';
import {PlayerSeekBar} from './PlayerSeekBar.js';
import {formatTime} from './format-time.js';
import {FullscreenIcon} from './icons.js';
import {useHoverState} from './use-hover-state.js';
import type {usePlayer} from './use-player.js';
import {
	X_PADDING,
	useVideoControlsResize,
} from './use-video-controls-resize.js';
import type {Size} from './utils/use-element-size.js';

export type RenderPlayPauseButton = (props: {
	playing: boolean;
	isBuffering: boolean;
}) => ReactNode | null;

export type RenderFullscreenButton = (props: {
	isFullscreen: boolean;
}) => ReactNode;

const gradientSteps = [
	0, 0.013, 0.049, 0.104, 0.175, 0.259, 0.352, 0.45, 0.55, 0.648, 0.741, 0.825,
	0.896, 0.951, 0.987,
];

const gradientOpacities = [
	0, 8.1, 15.5, 22.5, 29, 35.3, 41.2, 47.1, 52.9, 58.8, 64.7, 71, 77.5, 84.5,
	91.9,
];

const globalGradientOpacity = 1 / 0.7;

const containerStyle: React.CSSProperties = {
	boxSizing: 'border-box',
	position: 'absolute',
	bottom: 0,
	width: '100%',
	paddingTop: 40,
	paddingBottom: 10,
	backgroundImage: `linear-gradient(to bottom,${gradientSteps
		.map((g, i) => {
			return `hsla(0, 0%, 0%, ${g}) ${
				gradientOpacities[i] * globalGradientOpacity
			}%`;
		})
		.join(', ')}, hsl(0, 0%, 0%) 100%)`,
	backgroundSize: 'auto 145px',
	display: 'flex',
	paddingRight: X_PADDING,
	paddingLeft: X_PADDING,
	flexDirection: 'column',
	transition: 'opacity 0.3s',
};

const controlsRow: React.CSSProperties = {
	display: 'flex',
	flexDirection: 'row',
	width: '100%',
	alignItems: 'center',
	justifyContent: 'center',
	userSelect: 'none',
	WebkitUserSelect: 'none',
};

const leftPartStyle: React.CSSProperties = {
	display: 'flex',
	flexDirection: 'row',
	userSelect: 'none',
	WebkitUserSelect: 'none',
	alignItems: 'center',
};

const xSpacer: React.CSSProperties = {
	width: 12,
};

const ySpacer: React.CSSProperties = {
	height: 8,
};

const flex1: React.CSSProperties = {
	flex: 1,
};

const fullscreen: React.CSSProperties = {};

export const Controls: React.FC<{
	readonly fps: number;
	readonly durationInFrames: number;
	readonly showVolumeControls: boolean;
	readonly player: ReturnType<typeof usePlayer>;
	readonly onFullscreenButtonClick: MouseEventHandler<HTMLButtonElement>;
	readonly isFullscreen: boolean;
	readonly allowFullscreen: boolean;
	readonly onExitFullscreenButtonClick: MouseEventHandler<HTMLButtonElement>;
	readonly spaceKeyToPlayOrPause: boolean;
	readonly onSeekEnd: () => void;
	readonly onSeekStart: () => void;
	readonly inFrame: number | null;
	readonly outFrame: number | null;
	readonly initiallyShowControls: number | boolean;
	readonly canvasSize: Size | null;
	readonly renderPlayPauseButton: RenderPlayPauseButton | null;
	readonly renderFullscreenButton: RenderFullscreenButton | null;
	readonly alwaysShowControls: boolean;
	readonly showPlaybackRateControl: boolean | number[];
	readonly containerRef: React.RefObject<HTMLDivElement>;
	readonly buffering: boolean;
	readonly hideControlsWhenPointerDoesntMove: boolean | number;
	readonly onPointerDown:
		| ((ev: PointerEvent | SyntheticEvent) => void)
		| undefined;
	readonly onDoubleClick: MouseEventHandler<HTMLDivElement> | undefined;
}> = ({
	durationInFrames,
	isFullscreen,
	fps,
	player,
	showVolumeControls,
	onFullscreenButtonClick,
	allowFullscreen,
	onExitFullscreenButtonClick,
	spaceKeyToPlayOrPause,
	onSeekEnd,
	onSeekStart,
	inFrame,
	outFrame,
	initiallyShowControls,
	canvasSize,
	renderPlayPauseButton,
	renderFullscreenButton,
	alwaysShowControls,
	showPlaybackRateControl,
	containerRef,
	buffering,
	hideControlsWhenPointerDoesntMove,
	onPointerDown,
	onDoubleClick,
}) => {
	const playButtonRef = useRef<HTMLButtonElement | null>(null);
	const frame = Internals.Timeline.useTimelinePosition();
	const [supportsFullscreen, setSupportsFullscreen] = useState(false);
	const hovered = useHoverState(
		containerRef,
		hideControlsWhenPointerDoesntMove,
	);

	const {maxTimeLabelWidth, displayVerticalVolumeSlider} =
		useVideoControlsResize({
			allowFullscreen,
			playerWidth: canvasSize?.width ?? 0,
		});
	const [shouldShowInitially, setInitiallyShowControls] = useState<
		boolean | number
	>(() => {
		if (typeof initiallyShowControls === 'boolean') {
			return initiallyShowControls;
		}

		if (typeof initiallyShowControls === 'number') {
			if (initiallyShowControls % 1 !== 0) {
				throw new Error(
					'initiallyShowControls must be an integer or a boolean',
				);
			}

			if (Number.isNaN(initiallyShowControls)) {
				throw new Error('initiallyShowControls must not be NaN');
			}

			if (!Number.isFinite(initiallyShowControls)) {
				throw new Error('initiallyShowControls must be finite');
			}

			if (initiallyShowControls <= 0) {
				throw new Error('initiallyShowControls must be a positive integer');
			}

			return initiallyShowControls;
		}

		throw new TypeError('initiallyShowControls must be a number or a boolean');
	});

	const containerCss: React.CSSProperties = useMemo(() => {
		// Hide if playing and mouse outside
		const shouldShow =
			hovered || !player.playing || shouldShowInitially || alwaysShowControls;
		return {
			...containerStyle,
			opacity: Number(shouldShow),
		};
	}, [hovered, shouldShowInitially, player.playing, alwaysShowControls]);

	useEffect(() => {
		if (playButtonRef.current && spaceKeyToPlayOrPause) {
			// This switches focus to play button when player.playing flag changes
			playButtonRef.current.focus({
				preventScroll: true,
			});
		}
	}, [player.playing, spaceKeyToPlayOrPause]);

	useEffect(() => {
		// Must be handled client-side to avoid SSR hydration mismatch
		setSupportsFullscreen(
			(typeof document !== 'undefined' &&
				(document.fullscreenEnabled ||
					// @ts-expect-error Types not defined
					document.webkitFullscreenEnabled)) ??
				false,
		);
	}, []);

	useEffect(() => {
		if (shouldShowInitially === false) {
			return;
		}

		const time = shouldShowInitially === true ? 2000 : shouldShowInitially;
		const timeout = setTimeout(() => {
			setInitiallyShowControls(false);
		}, time);

		return () => {
			clearInterval(timeout);
		};
	}, [shouldShowInitially]);

	const timeLabel: React.CSSProperties = useMemo(() => {
		return {
			color: 'white',
			fontFamily: 'sans-serif',
			fontSize: 14,
			maxWidth: maxTimeLabelWidth === null ? undefined : maxTimeLabelWidth,
			overflow: 'hidden',
			textOverflow: 'ellipsis',
		};
	}, [maxTimeLabelWidth]);

	const playbackRates = useMemo(() => {
		if (showPlaybackRateControl === true) {
			return [0.5, 0.8, 1, 1.2, 1.5, 1.8, 2, 2.5, 3];
		}

		if (Array.isArray(showPlaybackRateControl)) {
			for (const rate of showPlaybackRateControl) {
				if (typeof rate !== 'number') {
					throw new Error(
						'Every item in showPlaybackRateControl must be a number',
					);
				}

				if (rate <= 0) {
					throw new Error(
						'Every item in showPlaybackRateControl must be positive',
					);
				}
			}

			return showPlaybackRateControl;
		}

		return null;
	}, [showPlaybackRateControl]);

	const ref = useRef<HTMLDivElement | null>(null);
	const flexRef = useRef<HTMLDivElement | null>(null);

	const onPointerDownIfContainer: React.PointerEventHandler<HTMLDivElement> =
		useCallback(
			(e) => {
				// Only if pressing the container
				if (e.target === ref.current || e.target === flexRef.current) {
					onPointerDown?.(e);
				}
			},
			[onPointerDown],
		);

	const onDoubleClickIfContainer: MouseEventHandler<HTMLDivElement> =
		useCallback(
			(e) => {
				// Only if pressing the container
				if (e.target === ref.current || e.target === flexRef.current) {
					onDoubleClick?.(e);
				}
			},
			[onDoubleClick],
		);

	return (
		<div
			ref={ref}
			style={containerCss}
			onPointerDown={onPointerDownIfContainer}
			onDoubleClick={onDoubleClickIfContainer}
			className='CONTROLSFRAME'
		>
			<div ref={flexRef} style={controlsRow}>
				<div style={leftPartStyle}>
					<button
						ref={playButtonRef}
						type="button"
						style={playerButtonStyle}
						onClick={player.playing ? player.pause : player.play}
						aria-label={player.playing ? 'Pause video' : 'Play video'}
						title={player.playing ? 'Pause video' : 'Play video'}
					>
						{renderPlayPauseButton === null ? (
							<DefaultPlayPauseButton
								buffering={buffering}
								playing={player.playing}
							/>
						) : (
							renderPlayPauseButton({
								playing: player.playing,
								isBuffering: buffering,
							}) ?? (
								<DefaultPlayPauseButton
									buffering={buffering}
									playing={player.playing}
								/>
							)
						)}
					</button>
					{showVolumeControls ? (
						<>
							<div style={xSpacer} />
							<MediaVolumeSlider
								displayVerticalVolumeSlider={displayVerticalVolumeSlider}
							/>
						</>
					) : null}
					<div style={xSpacer} />
					<div style={timeLabel}>
						{formatTime(frame / fps)} / {formatTime(durationInFrames / fps)}
					</div>
					<div style={xSpacer} />
				</div>
				<div style={flex1} />
				{playbackRates && canvasSize && (
					<PlaybackrateControl
						canvasSize={canvasSize}
						playbackRates={playbackRates}
					/>
				)}
				{playbackRates && supportsFullscreen && allowFullscreen ? (
					<div style={xSpacer} />
				) : null}
				<div style={fullscreen}>
					{supportsFullscreen && allowFullscreen ? (
						<button
							type="button"
							aria-label={isFullscreen ? 'Exit fullscreen' : 'Enter Fullscreen'}
							title={isFullscreen ? 'Exit fullscreen' : 'Enter Fullscreen'}
							style={playerButtonStyle}
							onClick={
								isFullscreen
									? onExitFullscreenButtonClick
									: onFullscreenButtonClick
							}
						>
							{renderFullscreenButton === null ? (
								<FullscreenIcon isFullscreen={isFullscreen} />
							) : (
								renderFullscreenButton({isFullscreen})
							)}
						</button>
					) : null}
				</div>
			</div>
			<div style={ySpacer} />
			<PlayerSeekBar
				onSeekEnd={onSeekEnd}
				onSeekStart={onSeekStart}
				durationInFrames={durationInFrames}
				inFrame={inFrame}
				outFrame={outFrame}
			/>
		</div>
	);
};
