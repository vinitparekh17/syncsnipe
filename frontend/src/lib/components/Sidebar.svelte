<script lang="ts">
	import type { NavItem } from '$lib/types/dashboard';
	import { theme } from '$lib/stores/theme';
	import Help from './icons/Help.svelte';
	import History from './icons/History.svelte';
	import Home from './icons/Home.svelte';
	import Network from './icons/Network.svelte';
	import Setting from './icons/Setting.svelte';
	import SyncProfile from './icons/SyncProfile.svelte';

	export let path: string;

	const navItems: NavItem[] = [
		{ href: '#stats',    text: 'Dashboard', icon: Home        },
		{ href: '#profiles', text: 'Profiles',  icon: SyncProfile },
		{ href: '#activity', text: 'Activity',  icon: History     },
		{ href: '#settings', text: 'Settings',  icon: Setting     },
		{ href: '#network',  text: 'Network',   icon: Network     },
		{ href: '#help',     text: 'Help',      icon: Help        }
	];

	$: isDark = $theme === 'dark';

	// Active item: Dashboard is active when at root with no hash
	function isActive(href: string): boolean {
		if (typeof window === 'undefined') return href === '#stats';
		const hash = window.location.hash;
		if (href === '#stats') return !hash || hash === '#stats';
		return hash === href;
	}

	// Track hash changes for nav highlighting
	let currentHash = '';
	if (typeof window !== 'undefined') {
		currentHash = window.location.hash;
		window.addEventListener('hashchange', () => {
			currentHash = window.location.hash;
		});
	}

	$: activeHash = currentHash;
	$: getActive = (href: string) => {
		if (href === '#stats') return !activeHash || activeHash === '#stats';
		return activeHash === href;
	};
</script>

<aside
	id="sidebar"
	class="fixed top-0 left-0 z-40 flex h-screen w-64 flex-col border-r"
	style="
		background: var(--sidebar-bg);
		border-right-color: var(--sidebar-border);
		backdrop-filter: var(--blur-material);
		-webkit-backdrop-filter: var(--blur-material);
		color: var(--text-primary);
	"
>
	<!-- Header / Logo -->
	<div
		class="flex items-center justify-between border-b px-6 py-5"
		style="border-bottom-color: var(--border-subtle);"
	>
		<div class="flex items-center space-x-3">
			<div class="rounded-xl bg-gradient-to-tr from-[#8A4FFF] to-[#20E6B5] p-2"
				style="box-shadow: 0 4px 14px rgba(138,79,255,0.30);">
				<svg class="h-6 w-6 text-white" viewBox="0 0 24 24" fill="none">
					<path d="M12 2L20 7V17L12 22L4 17V7L12 2Z"
						stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
					<path d="M12 7V12M12 17V12M7 9.5V14.5M17 9.5V14.5"
						stroke="currentColor" stroke-width="2" stroke-linecap="round" />
				</svg>
			</div>
			<div>
				<span class="text-lg font-bold tracking-tight bg-gradient-to-r from-[#8A4FFF] to-[#20E6B5] bg-clip-text text-transparent">
					SyncSnipe
				</span>
				<div class="text-[10px] tracking-widest font-semibold uppercase text-tertiary">
					Engine v0.1
				</div>
			</div>
		</div>
	</div>

	<!-- Navigation -->
	<nav class="mt-4 flex-1 px-3 space-y-0.5 overflow-y-auto">
		<ul>
			{#each navItems as navItem}
				{@const active = getActive(navItem.href)}
				<li>
					<a
						href={navItem.href}
						class="flex items-center space-x-3 rounded-xl px-4 py-2.5 text-sm font-medium border
							transition-all duration-150 btn-press"
						style="
							color: {active ? 'var(--text-primary)' : 'var(--text-secondary)'};
							background: {active ? 'var(--bg-inset-hover)' : 'transparent'};
							border-color: {active ? 'var(--border-default)' : 'transparent'};
						"
						on:mouseenter={(e) => {
							if (!active) {
								(e.currentTarget as HTMLElement).style.background = 'var(--bg-inset)';
								(e.currentTarget as HTMLElement).style.color = 'var(--text-primary)';
							}
						}}
						on:mouseleave={(e) => {
							if (!active) {
								(e.currentTarget as HTMLElement).style.background = 'transparent';
								(e.currentTarget as HTMLElement).style.color = 'var(--text-secondary)';
							}
						}}
					>
						<svelte:component this={navItem.icon} />
						<span class="tracking-tight">{navItem.text}</span>
					</a>
				</li>
			{/each}
		</ul>
	</nav>

	<!-- Footer -->
	<div class="p-4 border-t space-y-2" style="border-top-color: var(--border-subtle);">
		<!-- Theme toggle -->
		<button
			id="theme-toggle"
			on:click={theme.toggle}
			aria-label="Switch to {isDark ? 'light' : 'dark'} mode"
			title="Switch to {isDark ? 'light' : 'dark'} mode"
			class="w-full flex items-center justify-between rounded-xl border px-3 py-2.5 text-xs font-medium btn-press"
			style="
				background: var(--bg-inset);
				border-color: var(--border-subtle);
				color: var(--text-secondary);
			"
			on:mouseenter={(e) => {
				(e.currentTarget as HTMLElement).style.background = 'var(--bg-inset-hover)';
				(e.currentTarget as HTMLElement).style.color = 'var(--text-primary)';
			}}
			on:mouseleave={(e) => {
				(e.currentTarget as HTMLElement).style.background = 'var(--bg-inset)';
				(e.currentTarget as HTMLElement).style.color = 'var(--text-secondary)';
			}}
		>
			<span class="flex items-center space-x-2">
				{#if isDark}
					<svg class="h-4 w-4 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<circle cx="12" cy="12" r="4" />
						<path stroke-linecap="round" d="M12 2v2M12 20v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M2 12h2M20 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42" />
					</svg>
					<span>Light Mode</span>
				{:else}
					<svg class="h-4 w-4 text-[#8A4FFF]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z" />
					</svg>
					<span>Dark Mode</span>
				{/if}
			</span>
			<!-- Animated pill -->
			<div
				class="relative h-5 w-9 rounded-full border transition-colors duration-300"
				style="
					background: {isDark ? 'rgba(139,79,255,0.15)' : 'rgba(32,230,181,0.15)'};
					border-color: {isDark ? 'rgba(139,79,255,0.3)' : 'rgba(32,230,181,0.3)'};
				"
			>
				<div
					class="absolute top-0.5 h-4 w-4 rounded-full shadow-sm transition-all duration-300"
					style="
						background: {isDark ? '#8A4FFF' : '#20E6B5'};
						transform: translateX({isDark ? '16px' : '2px'});
						transition-timing-function: cubic-bezier(0.34,1.56,0.64,1);
					"
				></div>
			</div>
		</button>

		<!-- Status badge -->
		<div
			class="rounded-xl border p-3 text-xs text-center"
			style="background: var(--bg-inset); border-color: var(--border-subtle);"
		>
			<div class="font-medium text-primary">Local Mode Active</div>
			<div class="text-[11px] mt-0.5 text-tertiary">SQLite sync rules online</div>
		</div>
	</div>
</aside>
