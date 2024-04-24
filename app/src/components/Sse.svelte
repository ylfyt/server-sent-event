<script lang="ts">
	import { onMount } from "svelte";

	type SseData<T> = {
		Id: number;
		Path: string;
		Data: T;
	};

	let count: number | null = null;
	let sse: EventSource | undefined = undefined;
	onMount(() => {
		start();
	});

	const start = () => {
		sse = new EventSource("http://localhost:8080/sse");
		sse.addEventListener("error", (e) => {
			console.log("err", e);
		});
		sse.addEventListener("message", (e) => {
			try {
				const res = JSON.parse(e.data) as SseData<number>;
				count = res.Data;
			} catch (error) {
				console.error(error);
				count = null;
			}
		});
		sse.addEventListener("open", (e) => {
			console.log("open", e);
		});
	};

	const stop = () => {
		sse?.close();
		sse = undefined;
		count = null;
	};

	const inc = async () => {
		fetch("http://localhost:8080/inc", {
			method: "post",
		});
	};
</script>

<div>
	<div>Count is {count}</div>
	<br />
	<button on:click={inc}>Inc</button>
	<br />
	<br />
	{#if !sse}
		<button on:click={start}> Start </button>
	{:else}
		<button on:click={stop}> Stop </button>
	{/if}
</div>
