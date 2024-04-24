<script lang="ts">
	import { onMount } from "svelte";

	let message = "Please wait...";
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
			const data = e.data;
			const date = new Date(parseInt(data));
			message = date.toLocaleString();
		});
		sse.addEventListener("open", (e) => {
			console.log("open", e);
		});
	};

	const stop = () => {
		sse?.close();
		sse = undefined;
	};
</script>

<div>
	<div>{message}</div>
	<br />
	<div>
		{#if !sse}
			<button on:click={start}> Start </button>
		{:else}
			<button on:click={stop}> Stop </button>
		{/if}
	</div>
</div>
