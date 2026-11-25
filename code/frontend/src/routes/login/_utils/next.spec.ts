import { describe, expect, it } from 'vitest';
import { safeNextPath as parse } from './next';

const ORIGIN = 'https://gostrategy.test';
const safeNextPath = (raw: string | null | undefined) => parse(raw, ORIGIN);

describe('safeNextPath', () => {
	it.each([
		['/friends/add?code=K7QXM2PA', '/friends/add?code=K7QXM2PA'],
		['/profile', '/profile'],
		['/select-ai?mode=human_vs_ai#top', '/select-ai?mode=human_vs_ai#top']
	])('keeps same-origin path %s', (input, expected) => {
		expect(safeNextPath(input)).toBe(expected);
	});

	it.each([
		['protocol-relative', '//evil.com'],
		['backslash host', '/\\evil.com'],
		['tab-smuggled host', '/\t/evil.com'],
		['newline-smuggled host', '/\n/evil.com'],
		['triple slash', '///evil.com'],
		['dot segment to protocol-relative', '/.//evil.com'],
		['dot-dot segment to protocol-relative', '/..//evil.com'],
		['encoded dot segment to protocol-relative', '/%2e//evil.com'],
		['dot segment to triple slash', '/.///evil.com'],
		['absolute url', 'https://evil.com'],
		['javascript scheme', 'javascript:alert(1)'],
		['relative path', 'friends'],
		['empty string', ''],
		['unparseable host', '//[']
	])('falls back to / for %s', (_name, input) => {
		expect(safeNextPath(input)).toBe('/');
	});

	it.each([null, undefined])('falls back to / for %s', (input) => {
		expect(safeNextPath(input)).toBe('/');
	});

	it('rejects a host that merely starts with the origin', () => {
		expect(parse('//gostrategy.test.evil.com', ORIGIN)).toBe('/');
	});
});
