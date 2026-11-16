import { describe, expect, it } from 'vitest';
import { getPasswordStrengthScore, validatePasswordInput } from './authValidation';

describe('authValidation', () => {
	describe('validatePasswordInput', () => {
		it('fails passwords shorter than min length', () => {
			const checks = validatePasswordInput('Aa1!');
			expect(checks.isMinLength).toBe(false);
		});

		it('validates uppercase requirement', () => {
			expect(validatePasswordInput('lowercase1!').hasUppercase).toBe(false);
			expect(validatePasswordInput('Uppercase1!').hasUppercase).toBe(true);
		});

		it('validates lowercase requirement', () => {
			expect(validatePasswordInput('UPPERCASE1!').hasLowercase).toBe(false);
			expect(validatePasswordInput('UPPERcase1!').hasLowercase).toBe(true);
		});

		it('validates number requirement', () => {
			expect(validatePasswordInput('NoDigitsHere!').hasNumber).toBe(false);
			expect(validatePasswordInput('WithDigit1!').hasNumber).toBe(true);
		});
	});

	describe('getPasswordStrengthScore', () => {
		it('returns 0 for empty password', () => {
			expect(getPasswordStrengthScore('')).toBe(0);
		});

		it('returns maximum score of 4 for compliant strong password', () => {
			expect(getPasswordStrengthScore('StrongP@ss123')).toBe(4);
		});

		it('counts each satisfied criterion', () => {
			expect(getPasswordStrengthScore('abc')).toBe(1);
			expect(getPasswordStrengthScore('Abc1')).toBe(3);
		});
	});
});
