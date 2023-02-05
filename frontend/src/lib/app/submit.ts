import { get } from 'svelte/store';
import { answers, attemptTime, totalTime, userID } from '$lib/app/stores';

export async function submitToBackend() {
    return await fetch('/api/submit', {
        method: 'POST',
        headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json',
            'X-User-ID': get(userID),
            'X-Attempt-Time': get(attemptTime).toString(),
            'X-Total-Time': get(totalTime).toString(),
        },
        body: JSON.stringify(get(answers)),
    }).then(response => response.json());
}
