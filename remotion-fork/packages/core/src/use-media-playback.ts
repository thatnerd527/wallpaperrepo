import type {RefObject} from 'react';
import {useCallback, useContext, useEffect, useRef} from 'react';
import {useMediaStartsAt} from './audio/use-audio-frame.js';
import {useBufferUntilFirstFrame} from './buffer-until-first-frame.js';
import {BufferingContextReact} from './buffering.js';
import {playAndHandleNotAllowedError} from './play-and-handle-not-allowed-error.js';
import {
	TimelineContext,
	usePlayingState,
	useTimelinePosition,
} from './timeline-position-state.js';
import {useCurrentFrame} from './use-current-frame.js';
import {useMediaBuffering} from './use-media-buffering.js';
import {useRequestVideoCallbackTime} from './use-request-video-callback-time.js';
import {useVideoConfig} from './use-video-config.js';
import {getMediaTime} from './video/get-current-time.js';
import {isIosSafari} from './video/video-fragment.js';
import {warnAboutNonSeekableMedia} from './warn-about-non-seekable-media.js';

export const DEFAULT_ACCEPTABLE_TIMESHIFT = 0.45;

const seek = (
	mediaRef: RefObject<HTMLVideoElement | HTMLAudioElement>,
	time: number,
): void => {
	if (!mediaRef.current) {
		return;
	}

	// iOS seeking does not support multiple decimals
	const timeToSet = isIosSafari() ? Number(time.toFixed(1)) : time;
	mediaRef.current.currentTime = timeToSet;
};

export const useMediaPlayback = ({
	mediaRef,
	src,
	mediaType,
	playbackRate: localPlaybackRate,
	onlyWarnForMediaSeekingError,
	acceptableTimeshift,
	pauseWhenBuffering,
	isPremounting,
	debugSeeking,
	onAutoPlayError,
}: {
	mediaRef: RefObject<HTMLVideoElement | HTMLAudioElement>;
	src: string | undefined;
	mediaType: 'audio' | 'video';
	playbackRate: number;
	onlyWarnForMediaSeekingError: boolean;
	acceptableTimeshift: number;
	pauseWhenBuffering: boolean;
	isPremounting: boolean;
	debugSeeking: boolean;
	onAutoPlayError: null | (() => void);
}) => {
	const {playbackRate: globalPlaybackRate} = useContext(TimelineContext);
	const frame = useCurrentFrame();
	const absoluteFrame = useTimelinePosition();
	const [playing] = usePlayingState();
	const buffering = useContext(BufferingContextReact);
	const {fps} = useVideoConfig();
	const mediaStartsAt = useMediaStartsAt();
	const lastSeekDueToShift = useRef<number | null>(null);

	if (!buffering) {
		throw new Error(
			'useMediaPlayback must be used inside a <BufferingContext>',
		);
	}

	const currentTime = useRequestVideoCallbackTime(mediaRef, mediaType);

	const desiredUnclampedTime = getMediaTime({
		frame,
		playbackRate: localPlaybackRate,
		startFrom: -mediaStartsAt,
		fps,
	});

	const isMediaTagBuffering = useMediaBuffering({
		element: mediaRef,
		shouldBuffer: pauseWhenBuffering,
		isPremounting,
	});

	const isVariableFpsVideoMap = useRef<Record<string, boolean>>({});

	const onVariableFpsVideoDetected = useCallback(() => {
		if (!src) {
			return;
		}

		if (debugSeeking) {
			// eslint-disable-next-line no-console
			console.log(
				`Detected ${src} as a variable FPS video. Disabling buffering while seeking.`,
			);
		}

		isVariableFpsVideoMap.current[src] = true;
	}, [debugSeeking, src]);

	const {bufferUntilFirstFrame, isBuffering} = useBufferUntilFirstFrame({
		mediaRef,
		mediaType,
		onVariableFpsVideoDetected,
	});

	const playbackRate = localPlaybackRate * globalPlaybackRate;

	// For short audio, a lower acceptable time shift is used
	const acceptableTimeShiftButLessThanDuration = (() => {
		if (mediaRef.current?.duration) {
			return Math.min(
				mediaRef.current.duration,
				acceptableTimeshift ?? DEFAULT_ACCEPTABLE_TIMESHIFT,
			);
		}

		return acceptableTimeshift;
	})();

	useEffect(() => {
		if (!playing) {
			mediaRef.current?.pause();
			return;
		}

		const isPlayerBuffering = buffering.buffering.current;
		const isMediaTagBufferingOrStalled = isMediaTagBuffering || isBuffering();

		if (isPlayerBuffering && !isMediaTagBufferingOrStalled) {
			mediaRef.current?.pause();
		}
	}, [
		buffering.buffering,
		isBuffering,
		isMediaTagBuffering,
		mediaRef,
		playing,
	]);

	useEffect(() => {
		const tagName = mediaType === 'audio' ? '<Audio>' : '<Video>';
		if (!mediaRef.current) {
			throw new Error(`No ${mediaType} ref found`);
		}

		if (!src) {
			throw new Error(
				`No 'src' attribute was passed to the ${tagName} element.`,
			);
		}

		const playbackRateToSet = Math.max(0, playbackRate);
		if (mediaRef.current.playbackRate !== playbackRateToSet) {
			mediaRef.current.playbackRate = playbackRateToSet;
		}

		const {duration} = mediaRef.current;
		const shouldBeTime =
			!Number.isNaN(duration) && Number.isFinite(duration)
				? Math.min(duration, desiredUnclampedTime)
				: desiredUnclampedTime;

		const mediaTagTime = mediaRef.current.currentTime;
		const rvcTime = currentTime.current ?? null;

		const isVariableFpsVideo = isVariableFpsVideoMap.current[src];

		const timeShiftMediaTag = Math.abs(shouldBeTime - mediaTagTime);
		const timeShiftRvcTag = rvcTime ? Math.abs(shouldBeTime - rvcTime) : null;
		const timeShift =
			timeShiftRvcTag && !isVariableFpsVideo
				? timeShiftRvcTag
				: timeShiftMediaTag;

		if (debugSeeking) {
			// eslint-disable-next-line no-console
			console.log({
				mediaTagTime,
				rvcTime,
				shouldBeTime,
				state: mediaRef.current.readyState,
				playing: !mediaRef.current.paused,
				isVariableFpsVideo,
			});
		}

		if (
			timeShift > acceptableTimeShiftButLessThanDuration &&
			lastSeekDueToShift.current !== shouldBeTime
		) {
			// If scrubbing around, adjust timing
			// or if time shift is bigger than 0.45sec

			if (debugSeeking) {
				// eslint-disable-next-line no-console
				console.log('Seeking', {
					shouldBeTime,
					isTime: mediaTagTime,
					rvcTime,
					timeShift,
				});
			}

			seek(mediaRef, shouldBeTime);
			lastSeekDueToShift.current = shouldBeTime;
			if (playing && !isVariableFpsVideo) {
				bufferUntilFirstFrame(shouldBeTime);
				if (mediaRef.current.paused) {
					playAndHandleNotAllowedError(mediaRef, mediaType, onAutoPlayError);
				}
			}

			if (!onlyWarnForMediaSeekingError) {
				warnAboutNonSeekableMedia(
					mediaRef.current,
					onlyWarnForMediaSeekingError ? 'console-warning' : 'console-error',
				);
			}

			return;
		}

		const seekThreshold = playing ? 0.15 : 0.00001;

		// Only perform a seek if the time is not already the same.
		// Chrome rounds to 6 digits, so 0.033333333 -> 0.033333,
		// therefore a threshold is allowed.
		// Refer to the https://github.com/remotion-dev/video-buffering-example
		// which is fixed by only seeking conditionally.
		const makesSenseToSeek =
			Math.abs(mediaRef.current.currentTime - shouldBeTime) > seekThreshold;

		const isMediaTagBufferingOrStalled = isMediaTagBuffering || isBuffering();
		const isSomethingElseBuffering =
			buffering.buffering.current && !isMediaTagBufferingOrStalled;

		if (!playing || isSomethingElseBuffering) {
			if (makesSenseToSeek) {
				seek(mediaRef, shouldBeTime);
			}

			return;
		}

		// We assured we are in playing state
		if (
			(mediaRef.current.paused && !mediaRef.current.ended) ||
			absoluteFrame === 0
		) {
			if (makesSenseToSeek) {
				seek(mediaRef, shouldBeTime);
			}

			playAndHandleNotAllowedError(mediaRef, mediaType, onAutoPlayError);
			if (!isVariableFpsVideo) {
				bufferUntilFirstFrame(shouldBeTime);
			}
		}
	}, [
		absoluteFrame,
		acceptableTimeShiftButLessThanDuration,
		bufferUntilFirstFrame,
		buffering.buffering,
		currentTime,
		debugSeeking,
		desiredUnclampedTime,
		isBuffering,
		isMediaTagBuffering,
		mediaRef,
		mediaType,
		onlyWarnForMediaSeekingError,
		playbackRate,
		playing,
		src,
		onAutoPlayError,
	]);
};
