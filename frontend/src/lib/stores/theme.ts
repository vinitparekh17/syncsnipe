import { writable } from 'svelte/store';
import { browser } from '$app/environment';

function createThemeStore() {
	// Initial value: dark (SSR-safe)
	const { subscribe, set, update } = writable<'light' | 'dark'>('dark');

	function init() {
		if (!browser) return;
		const saved = localStorage.getItem('syncsnipe-theme') as 'light' | 'dark' | null;
		const preferred = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
		const resolved = saved ?? preferred;
		set(resolved);
		document.documentElement.setAttribute('data-theme', resolved);
	}

	function toggle() {
		update((current) => {
			const next = current === 'dark' ? 'light' : 'dark';
			if (browser) {
				document.documentElement.setAttribute('data-theme', next);
				localStorage.setItem('syncsnipe-theme', next);
			}
			return next;
		});
	}

	return { subscribe, toggle, init };
}

export const theme = createThemeStore();
