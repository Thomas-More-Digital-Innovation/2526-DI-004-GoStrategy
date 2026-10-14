import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { readFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';

declare const process: {
	env: Record<string, string | undefined>;
};

const pkg = JSON.parse(
	readFileSync(fileURLToPath(new URL('package.json', import.meta.url)), 'utf8')
);

const commitSha = process.env.VITE_COMMIT_SHA ? process.env.VITE_COMMIT_SHA.slice(0, 7) : 'local';

const repoName =
	process.env.VITE_REPO_NAME || 'Thomas-More-Digital-Innovation/2526-DI-004-Stratego';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	define: {
		__VERSION__: JSON.stringify(pkg.version),
		__COMMIT_SHA__: JSON.stringify(commitSha),
		__REPO_NAME__: JSON.stringify(repoName)
	}
});
