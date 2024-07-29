import {useEffect, useRef} from 'react';

const getIsBackgrounded = () => {
	if (typeof document === 'undefined') {
		return false;
	}

	return document.visibilityState === 'hidden';
};

export const useIsBackgrounded = () => {
	const isBackgrounded = useRef(getIsBackgrounded());

	useEffect(() => {
		const onVisibilityChange = () => {
			isBackgrounded.current = getIsBackgrounded();
		};

		document.addEventListener('visibilitychange', onVisibilityChange);

		return () => {
			document.removeEventListener('visibilitychange', onVisibilityChange);
		};
	}, []);

	return isBackgrounded;
};
