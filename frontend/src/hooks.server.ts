import type { HandleFetch } from '@sveltejs/kit';
import { env } from '$env/dynamic/public';

export const handleFetch: HandleFetch = async ({ request, fetch, event }) => {
    const regex = /^(http:..sveltekit-prerender)?\/api\//;
    const apiBackend = env.PUBLIC_BACKEND;

    if (apiBackend && regex.test(request.url)) {
        request = new Request(
            apiBackend + request.url.replace(regex, '/api/'),
            request
        );
        request.headers.set('origin', event.url.origin);
    }

    return fetch(request);
};
