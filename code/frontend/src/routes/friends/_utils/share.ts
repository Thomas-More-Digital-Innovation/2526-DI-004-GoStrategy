const TITLE = 'GoStrategy';

/** The slice of `navigator` the share flow needs, so tests can pass a fake. */
export interface ShareTarget {
	share?: (data: ShareData) => Promise<void>;
	canShare?: (data: ShareData) => boolean;
	clipboard: { writeText: (text: string) => Promise<void> };
}

export type ShareOutcome = 'shared' | 'copied' | 'cancelled';

export function buildInviteUrl(origin: string, code: string): string {
	return `${origin}/friends/add?code=${encodeURIComponent(code)}`;
}

// the url stays out of `text`: some share targets print both and show the link twice
export function buildShareText(code: string): string {
	return `Play GoStrategy with me, add me as a friend! Friend code: ${code}`;
}

/**
 * Opens the native share sheet when the browser offers one, else copies text and link.
 * Closing the sheet is a normal outcome, any other share failure falls back to the clipboard.
 */
export async function shareInvite(
	target: ShareTarget,
	invite: { code: string; origin: string }
): Promise<ShareOutcome> {
	const text = buildShareText(invite.code);
	const url = buildInviteUrl(invite.origin, invite.code);
	const data: ShareData = { title: TITLE, text, url };

	if (target.share && target.canShare?.(data)) {
		try {
			await target.share(data);
			return 'shared';
		} catch (e) {
			if (e instanceof DOMException && e.name === 'AbortError') return 'cancelled';
		}
	}
	await target.clipboard.writeText(`${text} ${url}`);
	return 'copied';
}
