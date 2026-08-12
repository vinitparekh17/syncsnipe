<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import AlertTriangle from '$lib/components/icons/AlertTriangle.svelte';
	import File from '$lib/components/icons/File.svelte';
	import Panding from '$lib/components/icons/Panding.svelte';
	import Storage from '$lib/components/icons/Storage.svelte';
	import Syncing from '$lib/components/icons/Syncing.svelte';
	import SyncProfile from '$lib/components/icons/SyncProfile.svelte';
	import ProfileCard from '$lib/components/profile-card/ProfileCard.svelte';
	import ActivityItem from '$lib/components/recent-activity/ActivityItem.svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import StateCard from '$lib/components/state-card/StateCard.svelte';
	import CreateProfileModal from '$lib/components/modals/CreateProfileModal.svelte';
	import AddSyncRuleModal from '$lib/components/modals/AddSyncRuleModal.svelte';

	import {
		getDashboardStats,
		getProfiles,
		getRecentActivities,
		formatBytes,
		type DashboardStats,
		type ProfileWithRules,
		type ActivityItem as ActivityType
	} from '$lib/api';

	let stats: DashboardStats | null = null;
	let profiles: ProfileWithRules[] = [];
	let activities: ActivityType[] = [];
	let loading = true;
	let error = '';

	let isCreateProfileOpen = false;
	let isAddRuleOpen = false;

	async function loadData() {
		loading = true;
		error = '';
		try {
			const [s, p, a] = await Promise.all([
				getDashboardStats(),
				getProfiles(),
				getRecentActivities()
			]);
			stats = s;
			profiles = p;
			activities = a;
		} catch (e: any) {
			error = e.message || 'Failed to connect to SyncSnipe engine';
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		loadData();
	});

	$: stateCards = stats
		? [
				{ label: 'Total Files',     value: stats.totalFiles,                        icon: File          },
				{ label: 'Active Profiles', value: stats.totalProfiles,                     icon: SyncProfile   },
				{ label: 'Conflicts',       value: stats.unresolvedConflicts,               icon: AlertTriangle },
				{ label: 'Active Rules',    value: stats.activeRules,                       icon: Panding       },
				{ label: 'Pending Syncs',   value: stats.pendingSyncs,                      icon: Syncing       },
				{ label: 'Storage Used',    value: formatBytes(stats.totalStorageBytes),     icon: Storage       }
		  ]
		: [];
</script>

<div id="root" class="min-h-screen" style="background: var(--bg-base); color: var(--text-primary);">
	<Sidebar path={page.url.pathname} />

	<main id="dashboard" class="ml-64 min-h-screen px-8 py-8">
		<div class="max-w-7xl mx-auto space-y-8">

			<!-- ── Header ────────────────────────────────────────── -->
			<div class="flex items-center justify-between">
				<div>
					<h1
						class="text-3xl font-bold bg-gradient-to-r from-[#8A4FFF] to-[#20E6B5] bg-clip-text text-transparent"
						style="letter-spacing: -0.02em; line-height: 1.1;"
					>
						Dashboard
					</h1>
					<p class="mt-1 text-sm text-secondary">
						Monitor and control your file synchronization rules in real time
					</p>
				</div>

				<div class="flex items-center space-x-3">
					<!-- New Profile button -->
					<button
						id="btn-new-profile"
						on:click={() => (isCreateProfileOpen = true)}
						class="flex items-center space-x-2 rounded-xl bg-gradient-to-r from-[#8A4FFF] to-[#20E6B5]
							px-4 py-2.5 text-xs font-semibold text-white btn-press hover:opacity-90"
						style="box-shadow: var(--shadow-brand);"
					>
						<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
						</svg>
						<span>New Profile</span>
					</button>

					<!-- Add Sync Rule button -->
					<button
						id="btn-add-rule"
						on:click={() => (isAddRuleOpen = true)}
						disabled={profiles.length === 0}
						class="flex items-center space-x-2 rounded-xl border border-themed bg-inset px-4 py-2.5
							text-xs font-semibold text-primary backdrop-blur-md btn-press
							hover:bg-inset-hover disabled:opacity-40"
					>
						<svg class="h-4 w-4 text-[#20E6B5]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
						</svg>
						<span>Add Sync Rule</span>
					</button>

					<!-- Refresh button -->
					<button
						id="btn-refresh"
						on:click={loadData}
						class="flex items-center rounded-xl border border-themed bg-inset p-2.5
							text-secondary backdrop-blur-md btn-press
							hover:bg-inset-hover hover:text-primary"
						title="Refresh Dashboard"
					>
						<svg class="h-4 w-4 {loading ? 'animate-spin' : ''}" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
								d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
						</svg>
					</button>
				</div>
			</div>

			<!-- ── Error banner ───────────────────────────────────── -->
			{#if error}
				<div class="rounded-2xl border border-red-500/20 bg-red-500/10 p-4 text-sm text-red-500">
					{error}
				</div>
			{/if}

			<!-- ── Stats Grid ─────────────────────────────────────── -->
			<!-- id="stats" allows sidebar "Dashboard" link to scroll here -->
			<div id="stats" class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{#if loading && !stats}
					{#each Array(6) as _}
						<div class="h-28 rounded-2xl border border-themed bg-inset animate-pulse"></div>
					{/each}
				{:else}
					{#each stateCards as stateCard}
						<StateCard {stateCard} />
					{/each}
				{/if}
			</div>

			<!-- ── Sync Profiles ──────────────────────────────────── -->
			<!-- id="profiles" — sidebar "Profiles" link scrolls here -->
			<div id="profiles" class="space-y-4">
				<div class="flex items-center justify-between">
					<h2 class="text-xl font-bold text-primary" style="letter-spacing: -0.015em;">
						Sync Profiles
					</h2>
					<span class="text-xs text-tertiary">
						{profiles.length} total profile{profiles.length === 1 ? '' : 's'}
					</span>
				</div>

				{#if loading && profiles.length === 0}
					<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
						{#each Array(3) as _}
							<div class="h-44 rounded-2xl border border-themed bg-inset animate-pulse"></div>
						{/each}
					</div>
				{:else if profiles.length === 0}
					<div class="rounded-2xl border border-dashed border-themed p-12 text-center surface">
						<div class="mx-auto flex h-12 w-12 items-center justify-center rounded-2xl bg-inset text-[#8A4FFF]">
							<svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
							</svg>
						</div>
						<h3 class="mt-4 text-base font-semibold text-primary">No Profiles Created Yet</h3>
						<p class="mt-1 text-xs text-secondary">Get started by creating your first file sync profile</p>
						<button
							on:click={() => (isCreateProfileOpen = true)}
							class="mt-5 inline-flex items-center space-x-2 rounded-xl bg-gradient-to-r from-[#8A4FFF] to-[#20E6B5]
								px-5 py-2.5 text-xs font-semibold text-white shadow-lg btn-press hover:opacity-90"
						>
							Create Profile
						</button>
					</div>
				{:else}
					<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
						{#each profiles as profile}
							<ProfileCard {profile} on:refresh={loadData} />
						{/each}
					</div>
				{/if}
			</div>

			<!-- ── Activity + System Status ──────────────────────── -->
			<!-- id="activity" — sidebar "Activity" link scrolls here -->
			<div id="activity" class="grid grid-cols-1 gap-6 lg:grid-cols-2 pb-8">
				<!-- Recent Activity -->
				<div class="rounded-2xl p-6 surface">
					<h2 class="text-lg font-bold text-primary mb-4" style="letter-spacing: -0.015em;">
						Recent Activity
					</h2>
					{#if activities.length > 0}
						<div class="space-y-3">
							{#each activities as activity}
								<ActivityItem {activity} />
							{/each}
						</div>
					{:else}
						<div class="p-8 text-center text-xs text-tertiary">
							No recent file sync activity recorded.
						</div>
					{/if}
				</div>

				<!-- System Status -->
				<div class="rounded-2xl p-6 surface">
					<h2 class="text-lg font-bold text-primary mb-4" style="letter-spacing: -0.015em;">
						System Status
					</h2>
					<div class="space-y-3">
						<div class="flex items-center justify-between rounded-xl border border-themed bg-inset p-3.5">
							<span class="text-sm font-medium text-primary">SQLite Database</span>
							<span class="rounded-full bg-emerald-500/15 px-2.5 py-0.5 text-[11px] font-semibold
								text-emerald-600 border border-emerald-500/25">Connected</span>
						</div>
						<div class="flex items-center justify-between rounded-xl border border-themed bg-inset p-3.5">
							<span class="text-sm font-medium text-primary">REST API Service</span>
							<span class="rounded-full bg-emerald-500/15 px-2.5 py-0.5 text-[11px] font-semibold
								text-emerald-600 border border-emerald-500/25">Online</span>
						</div>
						<div class="flex items-center justify-between rounded-xl border border-themed bg-inset p-3.5">
							<span class="text-sm font-medium text-primary">Initial Sync Mode</span>
							<span class="text-xs font-mono text-secondary">Mirror / Events Only</span>
						</div>
					</div>
				</div>
			</div>

		</div>
	</main>
</div>

<CreateProfileModal bind:isOpen={isCreateProfileOpen} on:created={loadData} />
<AddSyncRuleModal bind:isOpen={isAddRuleOpen} {profiles} on:added={loadData} />
