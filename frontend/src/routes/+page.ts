import { env } from '$env/dynamic/public';
import type { PageLoad } from './$types';
import PocketBase from 'pocketbase';

const pb = new PocketBase(env.PUBLIC_BACKEND);

export const load = (async () => {
    const settings = pb.collection('settings');
    let ret: { [key: string]: any } = {};

    for (let record of await settings.getFullList()) {
        let val: string | string[] = record.value;
        // for list of strings, join them into one HTML string
        ret[record.key] = typeof val == 'string' ? val : val.join('\n');
    }

    ret.questions = pb
        .collection('questions')
        .getFullList(undefined, { sort: 'index' })
        .then(list =>
            list.map(q => {
                return {
                    id: q.id,
                    question: q.question,
                    choices: q.choices,
                };
            })
        );

    return ret;
}) satisfies PageLoad;
