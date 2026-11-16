import { defineConfig, mergeConfig } from 'vitest/config';
import { playwright } from '@vitest/browser-playwright';
import viteConfig from './vite.config.ts';

// why: vitest ignores vite.config.ts once this file exists, so merge to keep plugins and defines
export default mergeConfig(
	viteConfig,
	defineConfig({
		test: {
			expect: { requireAssertions: true },
			coverage: {
				provider: 'v8',
				include: ['src/**/*.{ts,svelte}'],
				exclude: [
					'src/**/*.spec.ts',
					'src/**/fake-*.ts',
					'src/**/*.d.ts',
					'src/lib/types/**',
					'src/lib/data/**',
					'src/routes/**/+*.ts'
				]
			},
			projects: [
				{
					extends: true,
					test: {
						name: 'server',
						environment: 'node',
						include: ['src/**/*.spec.ts'],
						exclude: ['src/**/*.svelte.spec.ts']
					}
				},
				{
					extends: true,
					test: {
						name: 'client',
						include: ['src/**/*.svelte.spec.ts'],
						browser: {
							enabled: true,
							headless: true,
							provider: playwright(),
							instances: [{ browser: 'chromium' }]
						}
					}
				}
			]
		}
	})
);
