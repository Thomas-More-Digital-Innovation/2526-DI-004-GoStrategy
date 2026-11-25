const FALLBACK = '/';

/**
 * Returns a same-origin path+query+hash safe to pass to `goto`, or `/`.
 * Parsing with `URL` makes the browser's own normalisation decide:
 * `/\evil.com` and `/<tab>/evil.com` both collapse to `//evil.com` and leave the origin.
 */
export function safeNextPath(raw: string | null | undefined, origin: string): string {
	if (!raw?.startsWith('/')) return FALLBACK;
	try {
		const url = new URL(raw, origin);
		if (url.origin !== origin) return FALLBACK;
		const path = url.pathname + url.search + url.hash;
		return path.startsWith('//') ? FALLBACK : path;
	} catch {
		return FALLBACK;
	}
}
