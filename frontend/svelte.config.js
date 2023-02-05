import adapter from '@sveltejs/adapter-static';
import preprocess from 'svelte-preprocess';

const config = {
    preprocess: preprocess(),

    kit: {
        alias: {
            $lib: 'src/lib',
        },
        adapter: adapter({
            prerender: { entries: [] },
            fallback: 'index.html',
        }),
    },
};

export default config;
