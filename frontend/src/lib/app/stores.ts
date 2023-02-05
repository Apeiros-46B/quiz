import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { nanoid } from 'nanoid';

// make a writable store from a value stored in Local Storage
function locallyStored(key: string, defaultValue: any) {
    let val = defaultValue;

    if (browser) {
        let stored = localStorage.getItem(key);
        if (stored != undefined) {
            val = JSON.parse(stored);
        }
    }

    const store = writable(val);

    store.subscribe(n => {
        if (browser) {
            localStorage.setItem(key, JSON.stringify(n));
        }
    });

    return store;
}

// currently selected choices
let obj: { [key: string]: number } = {};
export const answers = writable(obj);
export const active = writable(-1);

// attempt state
export const attemptTime = writable(0);
export const attemptCount = locallyStored('attempts', 0);
export const totalTime = locallyStored('totalTime', 0);

// user state
export const userID = locallyStored('userID', nanoid());
