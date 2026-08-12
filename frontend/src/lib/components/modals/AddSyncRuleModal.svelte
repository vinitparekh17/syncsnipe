<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { addSyncRule, type ProfileWithRules } from '$lib/api';

	export let isOpen = false;
	export let profiles: ProfileWithRules[] = [];

	const dispatch = createEventDispatcher();
	let selectedProfile = '';
	let sourceDir = '';
	let targetDir = '';
	let initialSyncMode = 'mirror';
	let loading = false;
	let error = '';

	$: if (isOpen && profiles.length > 0 && !selectedProfile) {
		selectedProfile = profiles[0].name;
	}

	async function handleSubmit() {
		error = '';
		if (!selectedProfile)   { error = 'Please select a profile';       return; }
		if (!sourceDir.trim())  { error = 'Source directory is required';  return; }
		if (!targetDir.trim())  { error = 'Target directory is required';  return; }

		loading = true;
		try {
			await addSyncRule({ profile: selectedProfile, sourceDir: sourceDir.trim(), targetDir: targetDir.trim(), initialSyncMode });
			sourceDir = '';
			targetDir = '';
			dispatch('added');
			isOpen = false;
		} catch (e: any) {
			error = e.message || 'Failed to add sync rule';
		} finally {
			loading = false;
		}
	}

	function close() {
		if (loading) return;
		error = '';
		sourceDir = '';
		targetDir = '';
		isOpen = false;
	}
</script>

{#if isOpen}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center p-4"
		style="background: var(--bg-scrim); backdrop-filter: var(--blur-overlay); -webkit-backdrop-filter: var(--blur-overlay);"
		on:click|self={close}
		on:keydown={(e) => e.key === 'Escape' && close()}
		role="dialog"
		aria-modal="true"
		aria-labelledby="modal-addrule-title"
		tabindex="-1"
	>
		<div class="w-full max-w-lg overflow-hidden rounded-2xl border border-themed p-6 surface-raised">
			<!-- Header -->
			<div class="flex items-center justify-between pb-4 border-b border-subtle">
				<h3 id="modal-addrule-title" class="text-lg font-semibold text-primary" style="letter-spacing: -0.015em;">
					Add Sync Rule
				</h3>
				<button
					on:click={close}
					class="rounded-full p-1.5 text-secondary hover:bg-inset hover:text-primary transition-colors btn-press"
					aria-label="Close modal"
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
					</svg>
				</button>
			</div>

			<form on:submit|preventDefault={handleSubmit} class="mt-4 space-y-4">
				{#if error}
					<div class="rounded-lg border border-red-500/20 bg-red-500/10 p-3 text-sm text-red-500">
						{error}
					</div>
				{/if}

				<!-- Profile selector -->
				<div>
					<label for="profileSelect" class="block text-xs font-semibold tracking-widest uppercase text-tertiary mb-1.5">
						Target Profile
					</label>
					<select
						id="profileSelect"
						bind:value={selectedProfile}
						disabled={loading || profiles.length === 0}
						class="w-full rounded-xl border border-themed bg-inset px-4 py-2.5 text-sm text-primary
							focus:border-[#8A4FFF] focus:outline-none focus:ring-2 focus:ring-[rgba(138,79,255,0.2)]
							transition-all duration-150 disabled:opacity-50 appearance-none cursor-pointer"
					>
						{#each profiles as p}
							<option value={p.name} style="background: var(--bg-elevated); color: var(--text-primary);">{p.name}</option>
						{/each}
					</select>
				</div>

				<!-- Source directory -->
				<div>
					<label for="sourceDir" class="block text-xs font-semibold tracking-widest uppercase text-tertiary mb-1.5">
						Source Directory
					</label>
					<input
						id="sourceDir"
						type="text"
						bind:value={sourceDir}
						placeholder="e.g. C:/Users/name/Documents"
						disabled={loading}
						class="w-full rounded-xl border border-themed bg-inset px-4 py-2.5 text-sm text-primary
							placeholder-[color:var(--text-tertiary)]
							focus:border-[#8A4FFF] focus:outline-none focus:ring-2 focus:ring-[rgba(138,79,255,0.2)]
							transition-all duration-150 disabled:opacity-50"
					/>
				</div>

				<!-- Target directory -->
				<div>
					<label for="targetDir" class="block text-xs font-semibold tracking-widest uppercase text-tertiary mb-1.5">
						Target Directory
					</label>
					<input
						id="targetDir"
						type="text"
						bind:value={targetDir}
						placeholder="e.g. D:/Backups/Documents"
						disabled={loading}
						class="w-full rounded-xl border border-themed bg-inset px-4 py-2.5 text-sm text-primary
							placeholder-[color:var(--text-tertiary)]
							focus:border-[#8A4FFF] focus:outline-none focus:ring-2 focus:ring-[rgba(138,79,255,0.2)]
							transition-all duration-150 disabled:opacity-50"
					/>
				</div>

				<!-- Sync mode -->
				<div>
					<p class="text-xs font-semibold tracking-widest uppercase text-tertiary mb-2">Initial Sync Mode</p>
					<div class="grid grid-cols-2 gap-3" id="syncMode">
						<button
							type="button"
							on:click={() => (initialSyncMode = 'mirror')}
							class="rounded-xl border px-4 py-2.5 text-xs font-medium text-left btn-press"
							style="
								border-color: {initialSyncMode === 'mirror' ? '#8A4FFF' : 'var(--border-default)'};
								background: {initialSyncMode === 'mirror' ? 'rgba(138,79,255,0.1)' : 'var(--bg-inset)'};
								color: {initialSyncMode === 'mirror' ? 'var(--text-primary)' : 'var(--text-secondary)'};
							"
						>
							<div class="font-semibold text-sm">Mirror</div>
							<div class="text-[11px] text-tertiary mt-0.5">Copy all existing files immediately</div>
						</button>
						<button
							type="button"
							on:click={() => (initialSyncMode = 'events_only')}
							class="rounded-xl border px-4 py-2.5 text-xs font-medium text-left btn-press"
							style="
								border-color: {initialSyncMode === 'events_only' ? '#20E6B5' : 'var(--border-default)'};
								background: {initialSyncMode === 'events_only' ? 'rgba(32,230,181,0.1)' : 'var(--bg-inset)'};
								color: {initialSyncMode === 'events_only' ? 'var(--text-primary)' : 'var(--text-secondary)'};
							"
						>
							<div class="font-semibold text-sm">Events Only</div>
							<div class="text-[11px] text-tertiary mt-0.5">Only sync future changes</div>
						</button>
					</div>
				</div>

				<div class="flex items-center justify-end space-x-3 pt-2">
					<button
						type="button"
						on:click={close}
						disabled={loading}
						class="rounded-xl px-4 py-2 text-sm font-medium text-secondary hover:text-primary transition-colors btn-press"
					>
						Cancel
					</button>
					<button
						type="submit"
						disabled={loading}
						class="flex items-center rounded-xl bg-gradient-to-r from-[#8A4FFF] to-[#20E6B5]
							px-5 py-2 text-sm font-semibold text-white btn-press hover:opacity-90 disabled:opacity-50"
						style="box-shadow: var(--shadow-brand);"
					>
						{#if loading}
							<svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
							</svg>
						{/if}
						Add Sync Rule
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
