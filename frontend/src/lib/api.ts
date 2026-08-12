export type DashboardStats = {
	// API returns camelCase JSON keys from Go struct tags
	totalFiles: number;
	totalProfiles: number;
	activeRules: number;
	unresolvedConflicts: number;
	pendingSyncs: number;
	totalStorageBytes: number;
	// Also accept old snake_case keys for backwards compatibility
	total_files?: number;
	total_profiles?: number;
	active_rules?: number;
	unresolved_conflicts?: number;
	pending_syncs?: number;
	total_storage_bytes?: number;
};

export type SyncRule = {
	id: number;
	profileId: number;
	sourceDir: string;
	targetDir: string;
	status: number;
	initialSyncMode: string;
	lastRunSuccessful: boolean | null;
	createdAt: number;
	updatedAt: number;
};

export type ProfileWithRules = {
	id: number;
	name: string;
	createdAt: number;
	updatedAt: number;
	syncRules: SyncRule[];
};

export type ActivityItem = {
	label: string;
	details: string;
	color: string;
};

export async function getDashboardStats(): Promise<DashboardStats> {
	const res = await fetch('/api/dashboard/stats');
	if (!res.ok) throw new Error('Failed to fetch dashboard stats');
	return res.json();
}

export async function getProfiles(): Promise<ProfileWithRules[]> {
	const res = await fetch('/api/profiles');
	if (!res.ok) throw new Error('Failed to fetch profiles');
	return res.json();
}

export async function createProfile(name: string): Promise<void> {
	const res = await fetch('/api/profiles', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ name })
	});
	if (!res.ok) {
		const err = await res.json().catch(() => ({ error: 'Failed to create profile' }));
		throw new Error(err.error || 'Failed to create profile');
	}
}

export async function deleteProfile(name: string): Promise<void> {
	const res = await fetch(`/api/profiles?name=${encodeURIComponent(name)}`, {
		method: 'DELETE'
	});
	if (!res.ok) {
		const err = await res.json().catch(() => ({ error: 'Failed to delete profile' }));
		throw new Error(err.error || 'Failed to delete profile');
	}
}

export async function addSyncRule(rule: {
	profile: string;
	sourceDir: string;
	targetDir: string;
	initialSyncMode: string;
}): Promise<void> {
	const res = await fetch('/api/sync-rules', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(rule)
	});
	if (!res.ok) {
		const err = await res.json().catch(() => ({ error: 'Failed to add sync rule' }));
		throw new Error(err.error || 'Failed to add sync rule');
	}
}

export async function removeSyncRule(profile: string, sourceDir: string): Promise<void> {
	const res = await fetch('/api/sync-rules', {
		method: 'DELETE',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ profile, sourceDir })
	});
	if (!res.ok) {
		const err = await res.json().catch(() => ({ error: 'Failed to remove sync rule' }));
		throw new Error(err.error || 'Failed to remove sync rule');
	}
}

export async function getRecentActivities(): Promise<ActivityItem[]> {
	const res = await fetch('/api/activities');
	if (!res.ok) return [];
	return res.json();
}

export function formatBytes(bytes: number, decimals = 1): string {
	if (!bytes || bytes === 0) return '0 B';
	const k = 1024;
	const dm = decimals < 0 ? 0 : decimals;
	const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
	const i = Math.floor(Math.log(bytes) / Math.log(k));
	return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
}
