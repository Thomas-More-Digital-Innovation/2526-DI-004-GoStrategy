import { redirect } from '@sveltejs/kit';
import { resolve } from '$app/paths';
import type { PageLoad } from './$types';

// why: invite links already shared point here, the confirmation now lives in a modal on /friends
export const load: PageLoad = ({ url }) => {
	const code = url.searchParams.get('code');
	redirect(
		307,
		code ? `${resolve('friends')}?add=${encodeURIComponent(code)}` : resolve('friends')
	);
};
