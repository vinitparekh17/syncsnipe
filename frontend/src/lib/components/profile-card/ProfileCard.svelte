<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { ProfileWithRules } from '$lib/api';
	import { deleteProfile, removeSyncRule } from '$lib/api';

	export let profile: ProfileWithRules;

	const dispatch = createEventDispatcher();
	let deletingProfile = false;

	async function handleDeleteProfile() {
		if (!confirm(`Are you sure you want to delete profile "${profile.name}" and all its sync rules?`)) return;
		deletingProfile = true;
		try {
			await deleteProfile(profile.name);
			dispatch('refresh');
		} catch (e: any) {
			alert(e.message || 'Failed to delete profile');
		} finally {
			deletingProfile = false;
		}
	}

	async function handleRemoveRule(sourceDir: string) {
		if (!confirm(`Remove sync rule for "${sourceDir}"?`)) return;
		try {
			await removeSyncRule(profile.name, sourceDir);
			dispatch('refresh');
		} catch (e: any) {
			alert(e.message || 'Failed to remove sync rule');
		}
	}
</script>

<div class="group relative overflow-hidden rounded-2xl border border-themed p-5 surface card-lift">
	<!-- Header -->
	<div class="flex items-center justify-between pb-3 border-b border-subtle">
		<div class="flex items-center space-x-2.5">
			<div class="h-3 w-3 rounded-full bg-[#8A4FFF]" style="box-shadow: 0 0 8px rgba(138,79,255,0.5);"></div>
			<h3 class="text-base font-semibold text-primary" style="letter-spacing: -0.015em;">{profile.name}</h3>
		</div>
		<button
			on:click={handleDeleteProfile}
			disabled={deletingProfile}
			title="Delete Profile"
			class="rounded-lg p-1.5 text-tertiary opacity-0 group-hover:opacity-100
				hover:bg-red-500/15 hover:text-red-500 transition-all duration-150 btn-press"
		>
			<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
					d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
			</svg>
		</button>
	</div>

	<!-- Sync Rules -->
	<div class="mt-4 space-y-3">
		{#if profile.syncRules && profile.syncRules.length > 0}
			{#each profile.syncRules as rule}
				<div class="rounded-xl border border-subtle bg-inset p-3">
					<div class="flex items-center justify-between">
						<span
							class="inline-flex items-center rounded-md px-2 py-0.5 text-[10px] font-semibold tracking-wide uppercase
								{rule.initialSyncMode === 'mirror'
									? 'bg-[rgba(138,79,255,0.12)] text-[#8A4FFF] border border-[rgba(138,79,255,0.25)]'
									: 'bg-[rgba(32,230,181,0.12)] text-[#20E6B5] border border-[rgba(32,230,181,0.25)]'}"
						>
							{rule.initialSyncMode}
						</span>
						<button
							on:click={() => handleRemoveRule(rule.sourceDir)}
							class="text-xs text-tertiary hover:text-red-500 transition-colors duration-150 btn-press"
							title="Remove rule"
						>
							Remove
						</button>
					</div>

					<div class="mt-2 space-y-1 text-xs">
						<div class="flex items-center text-primary truncate" title={rule.sourceDir}>
							<svg class="mr-1.5 h-3.5 w-3.5 flex-shrink-0 text-[#8A4FFF]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
									d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
							</svg>
							<span class="truncate">{rule.sourceDir}</span>
						</div>
						<div class="flex items-center text-secondary truncate" title={rule.targetDir}>
							<svg class="mr-1.5 h-3.5 w-3.5 flex-shrink-0 text-[#20E6B5]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
									d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
							</svg>
							<span class="truncate">{rule.targetDir}</span>
						</div>
					</div>
				</div>
			{/each}
		{:else}
			<div class="rounded-xl border border-dashed border-subtle p-4 text-center text-xs text-tertiary">
				No sync rules configured yet.
			</div>
		{/if}
	</div>
</div>
