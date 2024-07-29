export const validateSingleFrame = (
	frame: unknown,
	variableName: string,
): number | null => {
	if (typeof frame === 'undefined' || frame === null) {
		return frame ?? null;
	}

	if (typeof frame !== 'number') {
		throw new TypeError(
			`"${variableName}" must be a number, but is ${JSON.stringify(frame)}`,
		);
	}

	if (Number.isNaN(frame)) {
		throw new TypeError(
			`"${variableName}" must not be NaN, but is ${JSON.stringify(frame)}`,
		);
	}

	if (!Number.isFinite(frame)) {
		throw new TypeError(
			`"${variableName}" must be finite, but is ${JSON.stringify(frame)}`,
		);
	}

	if (frame % 1 !== 0) {
		throw new TypeError(
			`"${variableName}" must be an integer, but is ${JSON.stringify(frame)}`,
		);
	}

	return frame;
};

export const validateInOutFrames = ({
	inFrame,
	durationInFrames,
	outFrame,
}: {
	inFrame: unknown;
	outFrame: unknown;
	durationInFrames: number;
}) => {
	const validatedInFrame = validateSingleFrame(inFrame, 'inFrame');
	const validatedOutFrame = validateSingleFrame(outFrame, 'outFrame');
	if (validatedInFrame === null && validatedOutFrame === null) {
		return;
	}

	// Must not be over the duration
	if (validatedInFrame !== null && validatedInFrame > durationInFrames - 1) {
		throw new Error(
			'inFrame must be less than (durationInFrames - 1), but is ' +
				validatedInFrame,
		);
	}

	if (validatedOutFrame !== null && validatedOutFrame > durationInFrames - 1) {
		throw new Error(
			'outFrame must be less than (durationInFrames - 1), but is ' +
				validatedOutFrame,
		);
	}

	// Must not be under 0
	if (validatedInFrame !== null && validatedInFrame < 0) {
		throw new Error(
			'inFrame must be greater than 0, but is ' + validatedInFrame,
		);
	}

	if (validatedOutFrame !== null && validatedOutFrame <= 0) {
		throw new Error(
			`outFrame must be greater than 0, but is ${validatedOutFrame}. If you want to render a single frame, use <Thumbnail /> instead.`,
		);
	}

	if (
		validatedOutFrame !== null &&
		validatedInFrame !== null &&
		validatedOutFrame <= validatedInFrame
	) {
		throw new Error(
			'outFrame must be greater than inFrame, but is ' +
				validatedOutFrame +
				' <= ' +
				validatedInFrame,
		);
	}
};
