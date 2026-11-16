import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { toastStore } from './toast.svelte';

describe('toastStore', () => {
	beforeEach(() => {
		vi.useFakeTimers();
		toastStore.toasts = [];
	});

	afterEach(() => {
		vi.useRealTimers();
	});

	it('adds a toast and removes it after its duration', () => {
		const id = toastStore.success('saved', 1000);
		expect(toastStore.toasts).toEqual([{ id, message: 'saved', type: 'success', duration: 1000 }]);
		vi.advanceTimersByTime(1000);
		expect(toastStore.toasts).toEqual([]);
	});

	it('keeps toasts with zero duration', () => {
		toastStore.info('sticky', 0);
		vi.runAllTimers();
		expect(toastStore.toasts).toHaveLength(1);
	});

	it('deduplicates identical messages', () => {
		toastStore.error('same');
		expect(toastStore.warning('same')).toBe('');
		expect(toastStore.toasts).toHaveLength(1);
	});

	it('removes by id', () => {
		const id = toastStore.add('x');
		toastStore.remove(id);
		expect(toastStore.toasts).toEqual([]);
	});

	describe('handleApiMessage', () => {
		it.each([
			['{"warning":"w"}', 'w', 'warning'],
			['{"info":"i"}', 'i', 'info'],
			['{"success":"s"}', 's', 'success'],
			['{"error":"e"}', 'e', 'error']
		])('maps %s to a %s toast', (input, message, type) => {
			toastStore.handleApiMessage(input);
			expect(toastStore.toasts[0]).toMatchObject({ message, type });
		});

		it('reads message and type from error objects', () => {
			const result = toastStore.handleApiMessage({ message: 'nope', type: 'warning' });
			expect(result).toMatchObject({ message: 'nope', type: 'warning' });
			expect(toastStore.toasts[0].type).toBe('warning');
		});

		it('stringifies unknown objects', () => {
			expect(toastStore.handleApiMessage({ code: 1 }).message).toBe('{"code":1}');
		});

		it('uses the fallback for empty input', () => {
			expect(toastStore.handleApiMessage({}, 'fallback').message).toBe('fallback');
			expect(toastStore.handleApiMessage(undefined).message).toBe('An error occurred');
		});
	});
});
