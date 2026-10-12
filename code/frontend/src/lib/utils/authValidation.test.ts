import { describe, expect, it } from 'vitest';
import {
	getPasswordStrengthScore,
	isPasswordStrong,
	validatePasswordInput
} from './authValidation';

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

		it('forbids leading and trailing spaces', () => {
			expect(validatePasswordInput(' Password123!').isValidFormat).toBe(false);
			expect(validatePasswordInput('Password123! ').isValidFormat).toBe(false);
			expect(validatePasswordInput('Pass word 123!').isValidFormat).toBe(true);
		});
	});

	describe('isPasswordStrong', () => {
		it('returns false for weak passwords', () => {
			expect(isPasswordStrong('short')).toBe(false);
			expect(isPasswordStrong('alllowercase123')).toBe(false);
			expect(isPasswordStrong('ALLUPPERCASE123')).toBe(false);
		});

		it('returns true when all criteria are met', () => {
			expect(isPasswordStrong('StrongP@ss123')).toBe(true);
		});
	});

	describe('getPasswordStrengthScore', () => {
		it('returns 0 for empty password', () => {
			expect(getPasswordStrengthScore('')).toBe(0);
		});

		it('returns maximum score of 4 for compliant strong password', () => {
			expect(getPasswordStrengthScore('StrongP@ss123')).toBe(4);
		});
	});
});
