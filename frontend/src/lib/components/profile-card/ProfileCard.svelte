<script lang="ts">
	import type { SyncProfileCardType } from '../../types/dashboard';
	export let profile: SyncProfileCardType;
</script>

<div
	class="animate__animated animate__fadeIn overflow-hidden rounded-lg bg-neutral-800 shadow-lg transition-shadow duration-300 hover:shadow-xl"
>
	<div class="p-4">
		<div class="mb-2 flex items-start justify-between">
			<h3 class="text-lg font-semibold">Project Backup</h3>
			<span class="rounded-full bg-green-500 px-2 py-1 text-xs font-medium"
				>{profile.status}</span
			>
		</div>
		<div class="mb-3">
			<div class="mb-1 flex items-center text-sm text-neutral-400">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="mr-2 h-4 w-4 text-[#8A4FFF]"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
					/>
				</svg>
				{profile.sourceDir}
			</div>
			<div class="flex items-center text-sm text-neutral-400">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="mr-2 h-4 w-4 text-[#20E6B5]"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"
					/>
				</svg>
				{profile.targetDir}
			</div>
		</div>
		<div class="mb-2">
			<div class="mb-1 flex justify-between text-xs">
				<span>Progress</span>
				<span>{profile.progress}%</span>
			</div>
			<div class="h-2 w-full rounded-full bg-neutral-700">
				<div
					class="h-2 rounded-full bg-gradient-to-r from-[#8A4FFF] to-[#20E6B5]"
					style="width: {profile.progress}%"
				></div>
			</div>
		</div>
		<div class="mb-4 text-xs text-neutral-400">
			Last synced: {profile.lastSync}
		</div>
		<div class="grid grid-cols-4 gap-2">
			<button class="rounded-md bg-neutral-700 p-2 hover:bg-neutral-600" aria-label="start-sync">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="mx-auto h-4 w-4"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"
					/>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
					/>
				</svg>
			</button>
			<button class="rounded-md bg-neutral-700 p-2 hover:bg-neutral-600" aria-label="stop-sync">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="mx-auto h-4 w-4"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z"
					/>
				</svg>
			</button>
			<button class="rounded-md bg-neutral-700 p-2 hover:bg-neutral-600" aria-label="edit-sync">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="mx-auto h-4 w-4"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
					/>
				</svg>
			</button>
			<button class="rounded-md bg-neutral-700 p-2 hover:bg-neutral-600" aria-label="delete-sync">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="mx-auto h-4 w-4"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
					/>
				</svg>
			</button>
		</div>
	</div>
</div>
