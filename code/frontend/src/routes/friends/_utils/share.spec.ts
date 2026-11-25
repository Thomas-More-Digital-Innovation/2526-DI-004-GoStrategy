import { describe, expect, it, vi } from 'vitest';
import { buildInviteUrl, buildShareText, shareInvite, type ShareTarget } from './share';

const invite = { code: 'K7QX-M2PA', origin: 'https://gostrategy.test' };
const url = 'https://gostrategy.test/friends/add?code=K7QX-M2PA';

function fakeTarget(overrides: Partial<ShareTarget> = {}) {
	const writeText = vi.fn().mockResolvedValue(undefined);
	const target: ShareTarget = { clipboard: { writeText }, ...overrides };
	return { target, writeText };
}

describe('share text and url', () => {
	it('keeps the url out of the text', () => {
		expect(buildShareText('K7QX-M2PA')).not.toContain('http');
	});

	it('builds a query-param invite url', () => {
		expect(buildInviteUrl('https://x.test', 'K7QX-M2PA')).toBe(
			'https://x.test/friends/add?code=K7QX-M2PA'
		);
	});
});

describe('shareInvite', () => {
	it('uses the native share sheet when the browser can share the data', async () => {
		const share = vi.fn().mockResolvedValue(undefined);
		const { target, writeText } = fakeTarget({ share, canShare: () => true });

		expect(await shareInvite(target, invite)).toBe('shared');
		expect(share).toHaveBeenCalledWith({
			title: 'GoStrategy',
			text: buildShareText('K7QX-M2PA'),
			url
		});
		expect(writeText).not.toHaveBeenCalled();
	});

	it('copies text and link when there is no share sheet', async () => {
		const { target, writeText } = fakeTarget();

		expect(await shareInvite(target, invite)).toBe('copied');
		expect(writeText).toHaveBeenCalledWith(`${buildShareText('K7QX-M2PA')} ${url}`);
	});

	it('copies when share exists but canShare says no', async () => {
		const share = vi.fn();
		const { target, writeText } = fakeTarget({ share, canShare: () => false });

		expect(await shareInvite(target, invite)).toBe('copied');
		expect(share).not.toHaveBeenCalled();
		expect(writeText).toHaveBeenCalled();
	});

	it('swallows a closed share sheet without copying', async () => {
		const share = vi.fn().mockRejectedValue(new DOMException('closed', 'AbortError'));
		const { target, writeText } = fakeTarget({ share, canShare: () => true });

		expect(await shareInvite(target, invite)).toBe('cancelled');
		expect(writeText).not.toHaveBeenCalled();
	});

	it('falls back to the clipboard when sharing fails for another reason', async () => {
		const share = vi.fn().mockRejectedValue(new DOMException('nope', 'NotAllowedError'));
		const { target, writeText } = fakeTarget({ share, canShare: () => true });

		expect(await shareInvite(target, invite)).toBe('copied');
		expect(writeText).toHaveBeenCalled();
	});

	it('lets a clipboard failure reach the caller', async () => {
		const { target, writeText } = fakeTarget();
		writeText.mockRejectedValue(new Error('denied'));

		await expect(shareInvite(target, invite)).rejects.toThrow('denied');
	});
});
