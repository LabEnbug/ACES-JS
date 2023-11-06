export const parseTime = (time, t) => {
	const date = new Date(time);
	const now = new Date();
	const diff = now.getTime() - date.getTime();
	const days = Math.floor(diff / (24 * 3600 * 1000));
	const hours = Math.floor((diff % (24 * 3600 * 1000)) / (3600 * 1000));
	const minutes = Math.floor((diff % (3600 * 1000)) / (60 * 1000));
	const seconds = Math.floor((diff % (60 * 1000)) / 1000);
	if (days > 0) {
		return `${days} ${t['time.daysAgo']}`;
	} else if (hours > 0) {
		return `${hours} ${t['time.hoursAgo']}`;
	} else if (minutes > 0) {
		return `${minutes} ${t['time.minutesAgo']}`;
	} else {
		return `${seconds} ${t['time.secondsAgo']}`;
	}
};