export const calculateNextFrame = ({
	time,
	currentFrame: startFrame,
	playbackSpeed,
	fps,
	actualLastFrame,
	actualFirstFrame,
	framesAdvanced,
	shouldLoop,
}: {
	time: number;
	currentFrame: number;
	playbackSpeed: number;
	fps: number;
	actualFirstFrame: number;
	actualLastFrame: number;
	framesAdvanced: number;
	shouldLoop: boolean;
}): {nextFrame: number; framesToAdvance: number; hasEnded: boolean} => {
	const op = playbackSpeed < 0 ? Math.ceil : Math.floor;
	const framesToAdvance =
		op((time * playbackSpeed) / (1000 / fps)) - framesAdvanced;

	const nextFrame = framesToAdvance + startFrame;
	const isCurrentFrameOutside =
		startFrame > actualLastFrame || startFrame < actualFirstFrame;
	const isNextFrameOutside =
		nextFrame > actualLastFrame || nextFrame < actualFirstFrame;

	const hasEnded = !shouldLoop && isNextFrameOutside && !isCurrentFrameOutside;
	if (playbackSpeed > 0) {
		// Play forwards
		if (isNextFrameOutside) {
			return {
				nextFrame: actualFirstFrame,
				framesToAdvance,
				hasEnded,
			};
		}

		return {nextFrame, framesToAdvance, hasEnded};
	}

	// Reverse playback
	if (isNextFrameOutside) {
		return {nextFrame: actualLastFrame, framesToAdvance, hasEnded};
	}

	return {nextFrame, framesToAdvance, hasEnded};
};
