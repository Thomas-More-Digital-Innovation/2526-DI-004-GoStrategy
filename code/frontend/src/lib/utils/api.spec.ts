import { describe, expect, it } from 'vitest';
import { parseApiMessage } from './api';

describe('parseApiMessage', () => {
	it('returns empty error for falsy input', () => {
		expect(parseApiMessage(null)).toEqual({ message: '', type: 'error' });
		expect(parseApiMessage('')).toEqual({ message: '', type: 'error' });
	});

	it('uses the first key of an object as the type', () => {
		expect(parseApiMessage({ warning: 'careful' })).toEqual({
			message: 'careful',
			type: 'warning'
		});
	});

	it('stringifies objects without a string message', () => {
		expect(parseApiMessage({ code: 42 })).toEqual({ message: '{"code":42}', type: 'error' });
	});

	it('returns empty message for empty object', () => {
		expect(parseApiMessage({})).toEqual({ message: '', type: 'error' });
	});

	it('parses JSON strings recursively', () => {
		expect(parseApiMessage('{"success":"saved"}')).toEqual({ message: 'saved', type: 'success' });
	});

	it('keeps plain text as an error message', () => {
		expect(parseApiMessage('boom')).toEqual({ message: 'boom', type: 'error' });
	});

	it('stringifies other primitives', () => {
		expect(parseApiMessage(404)).toEqual({ message: '404', type: 'error' });
	});
});
