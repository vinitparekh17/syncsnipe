<script lang="ts">
	import type { NavItem } from '$lib/types/dashboard';
	import Help from './icons/Help.svelte';
	import History from './icons/History.svelte';
	import Home from './icons/Home.svelte';
	import Network from './icons/Network.svelte';
	import Setting from './icons/Setting.svelte';
	import SyncProfile from './icons/SyncProfile.svelte';

	export let path: string;

	const navItems: NavItem[] = [
		{
			href: '/',
			text: 'Dashboard',
			icon: Home
		},
		{
			href: '#',
			text: 'Profiles',
			icon: SyncProfile
		},
		{
			href: '#',
			text: 'History',
			icon: History
		},
		{
			href: '#',
			text: 'Settings',
			icon: Setting
		},
		{
			href: '#',
			text: 'Network',
			icon: Network
		},
		{
			href: '#',
			text: 'Help/Support',
			icon: Help
		}
	];
</script>

<nav
	class="fixed top-0 left-0 flex h-screen w-64 flex-col bg-neutral-800 text-white shadow-lg dark:bg-neutral-800"
>
	<div class="flex items-center justify-between border-b border-neutral-700 p-4">
		<div class="flex items-center space-x-2">
			<svg class="h-8 w-8" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
				<path
					d="M12 2L20 7V17L12 22L4 17V7L12 2Z"
					stroke="url(#header-gradient)"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				/>
				<path
					d="M12 7V12M12 17V12M7 9.5V14.5M17 9.5V14.5"
					stroke="url(#header-gradient)"
					stroke-width="2"
					stroke-linecap="round"
				/>
				<defs>
					<linearGradient
						id="header-gradient"
						x1="0%"
						y1="0%"
						x2="100%"
						y2="0%"
						gradientUnits="userSpaceOnUse"
					>
						<stop offset="0%" stop-color="#8A4FFF" />
						<stop offset="100%" stop-color="#20E6B5" />
					</linearGradient>
				</defs>
			</svg>
			<span class="text-xl font-bold">SyncSnipe</span>
		</div>
		<!-- <button 
        type="button" 
        class="lg:hidden text-neutral-400 hover:text-white focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white"
        aria-controls="mobile-menu"
      >
        <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
        <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button> -->
	</div>

	<div class="mt-6 flex-1 overflow-y-auto">
		<ul>
			{#each navItems as navItem}
				<li
					class="mx-2 mb-1 rounded-lg px-4 py-3 transition-all duration-200 hover:bg-neutral-700 {navItem.href ==
						path && 'bg-neutral-700'}"
				>
					<a
						href={navItem.href}
						class="flex items-center space-x-3 text-neutral-300 hover:text-white"
					>
						<svelte:component this={navItem.icon} />
						<span>{navItem.text}</span>
					</a>
				</li>
			{/each}
		</ul>
	</div>

	<!-- Mobile menu -->
	<!-- x-transition:enter="transition ease-out duration-100 transform"
    x-transition:enter-start="opacity-0 scale-95"
    x-transition:enter-end="opacity-100 scale-100"
    x-transition:leave="transition ease-in duration-75 transform"
    x-transition:leave-start="opacity-100 scale-100"
    x-transition:leave-end="opacity-0 scale-95" -->
	<div id="mobile-menu" class="fixed inset-0 z-50 bg-neutral-800/80 backdrop-blur-lg lg:hidden">
		<div class="px-6 pt-20 pb-6">
			<ul class="space-y-2">
				{#each navItems as navItem}
					<li
						class="rounded-lg px-4 py-3 hover:bg-neutral-700 {navItem.href == path && 'bg-neutral-700'}"
					>
						<a
							href={navItem.href}
							class="flex items-center space-x-3 text-neutral-300 hover:text-white"
						>
							<svelte:component this={navItem.icon} />
							<span>{navItem.text}</span>
						</a>
					</li>
					
				{/each}
			</ul>
		</div>
	</div>
</nav>
