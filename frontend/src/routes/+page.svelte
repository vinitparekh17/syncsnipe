<script lang="ts">
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
	import type { RecentActivityItemType, StateCardType, SyncProfileCardType } from '$lib/types/dashboard';

	const mockArr: StateCardType[] = [
		{ label: 'Total Files', value: 5432, icon: File },
		{ label: 'Active Profiles', value: 3, icon: SyncProfile },
		{ label: 'Conflicts', value: 2, icon: AlertTriangle },
		{ label: 'Pending Syncs', value: 7, icon: Panding },
		{ label: 'Last Sync Time', value: 5, icon: Syncing }, // Could be replaced with a better icon
		{ label: 'Storage Used', value: 2.4, icon: Storage } // Consider formatting as "2.4 GB"
	];

	const mockProfileCards: SyncProfileCardType[] = [
		{
			sourceDir: '/home/user/documents',
			targetDir: '/backup/documents',
			status: 'active',
			progress: 85,
			lastSync: '2025-03-10T14:30:00Z'
		},
		{
			sourceDir: '/var/www/html',
			targetDir: '/backup/html',
			status: 'paused',
			progress: 40,
			lastSync: '2025-03-09T12:15:00Z'
		},
		{
			sourceDir: '/home/user/photos',
			targetDir: '/backup/photos',
			status: 'scheduled',
			progress: 10,
			lastSync: '2025-03-08T08:45:00Z'
		},
		{
			sourceDir: '/projects/code',
			targetDir: '/backup/code',
			status: 'ideal',
			progress: 0,
			lastSync: '2025-03-07T18:00:00Z'
		}
	];

	const mockActivities: RecentActivityItemType[] = [
		{
			label: "Project Backup sync completed",
			details: "Today, 2:30 PM • 1.2 GB transferred",
			color: "green"
		},
		{
			label: "Work Documents conflicts detected",
			details: "Today, 11:45 AM • 3 files in conflict",
			color: "red"
		},
		{
			label: "Music Library sync scheduled",
			details: "Today, 10:15 AM • Scheduled for 8:00 PM",
			color: "amber"
		},
		{
			label: "Personal Photos sync completed",
			details: "Yesterday, 8:15 AM • 3.5 GB transferred",
			color: "cyan"
		},
		{
			label: "Video Project sync failed",
			details: "Yesterday, 6:20 AM • Insufficient disk space",
			color: "violet"
		},
	];
</script>

<div id="root">
	<Sidebar path={page.url.pathname} />
	<section id="dashboard" class="ml-64 min-h-screen bg-neutral-900 text-white">
		<div class="container mx-auto px-4 py-8">
			<!-- Dashboard Header -->
			<div class="mb-8">
				<h1
					class="animate__animated animate__fadeIn mb-2 bg-gradient-to-r from-[#8A4FFF] to-[#20E6B5] bg-clip-text text-3xl font-bold text-transparent"
				>
					Dashboard
				</h1>
				<p class="text-neutral-400">Monitor and manage your file synchronization activities</p>
			</div>

			<!-- Quick Stats Widget -->
			<div
				class="animate__animated animate__fadeInUp mb-8 grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3"
			>
				{#each mockArr as stateCard}
					<StateCard {stateCard} />
				{/each}
			</div>

			<!-- Dashboard Actions -->
			<div class="animate__animated animate__fadeIn mb-8 flex flex-wrap gap-4">
				<button
					class="flex items-center rounded-md bg-gradient-to-r from-[#8A4FFF] to-[#20E6B5] px-4 py-2 text-white hover:opacity-90"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="mr-2 h-5 w-5"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 6v6m0 0v6m0-6h6m-6 0H6"
						/>
					</svg>
					Create New Profile
				</button>
				<button
					class="flex items-center rounded-md bg-neutral-800 px-4 py-2 text-white hover:bg-neutral-700"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="mr-2 h-5 w-5"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
						/>
					</svg>
					Sync All
				</button>
				<button
					class="flex items-center rounded-md bg-neutral-800 px-4 py-2 text-white hover:bg-neutral-700"
					id="togglePauseBtn"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="mr-2 h-5 w-5"
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
					Global Pause
				</button>
				<button
					class="flex items-center rounded-md bg-neutral-800 px-4 py-2 text-white hover:bg-neutral-700"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="mr-2 h-5 w-5"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
						/>
					</svg>
					Refresh Status
				</button>
			</div>

			<!-- Profile Cards -->
			<h2
				class="mb-4 bg-gradient-to-r from-[#8A4FFF] to-[#20E6B5] bg-clip-text text-xl font-semibold text-transparent"
			>
				Sync Profiles
			</h2>
			<div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
				<!-- Profile Card 1 -->
				{#each mockProfileCards as profile}
					<ProfileCard {profile} />
				{/each}
				<!-- <ProfileCard /> -->
			</div>

			<!-- Recent Activity -->
			<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
				<!-- Activity Timeline -->
				<div class="animate__animated animate__fadeIn rounded-lg bg-neutral-800 p-4 shadow-lg">
					<h2
						class="mb-4 bg-gradient-to-r from-[#8A4FFF] to-[#20E6B5] bg-clip-text text-xl font-semibold text-transparent"
					>
						Recent Activity
					</h2>
					<div class="space-y-4">
						{#each mockActivities as activity}
							<ActivityItem {activity} />
						{/each}
					</div>
					<button class="mt-4 w-full text-center text-sm text-[#20E6B5] hover:underline"
						>View Full History</button
					>
				</div>

				<!-- Sync Performance Chart -->
				<div class="animate__animated animate__fadeIn rounded-lg bg-neutral-800 p-4 shadow-lg">
					<h2
						class="mb-4 bg-gradient-to-r from-[#8A4FFF] to-[#20E6B5] bg-clip-text text-xl font-semibold text-transparent"
					>
						Sync Performance
					</h2>
					<div class="flex h-64 items-center justify-center" id="performanceChart">
						<canvas id="syncChart"></canvas>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Scripts -->
	<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
	<script>
		// Toggle Global Pause Button
		const togglePauseBtn = document.getElementById('togglePauseBtn');
		let isPaused = false;

		togglePauseBtn.addEventListener('click', function () {
			isPaused = !isPaused;
			if (isPaused) {
				this.innerHTML = `
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Resume All
            `;
				this.classList.remove('bg-neutral-800');
				this.classList.add('bg-[#8A4FFF]');
			} else {
				this.innerHTML = `
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Global Pause
            `;
				this.classList.remove('bg-[#8A4FFF]');
				this.classList.add('bg-neutral-800');
			}
		});

		// Performance Chart
		document.addEventListener('DOMContentLoaded', function () {
			const ctx = document.getElementById('syncChart').getContext('2d');

			// Chart data
			const syncData = {
				labels: ['1s', '2s', '3s', '4s', '5s', '6s', '7s'],
				datasets: [
					{
						label: 'Data Transferred (GB)',
						data: [3.2, 5.1, 2.8, 8.5, 4.2, 1.5, 6.3],
						borderColor: '#8A4FFF',
						backgroundColor: 'rgba(138, 79, 255, 0.1)',
						borderWidth: 2,
						tension: 0.4,
						fill: true
					}
				]
			};

			// Chart configuration
			const config = {
				type: 'line',
				data: syncData,
				options: {
					responsive: true,
					maintainAspectRatio: false,
					scales: {
						y: {
							beginAtZero: true,
							grid: {
								color: 'rgba(255, 255, 255, 0.1)'
							},
							ticks: {
								color: '#9ca3af'
							}
						},
						x: {
							grid: {
								color: 'rgba(255, 255, 255, 0.1)'
							},
							ticks: {
								color: '#9ca3af'
							}
						}
					},
					plugins: {
						legend: {
							labels: {
								color: '#ffffff'
							}
						}
					}
				}
			};

			// Create chart
			const syncChart = new Chart(ctx, config);
		});
	</script>
</div>
