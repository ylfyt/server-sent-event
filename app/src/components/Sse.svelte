<script lang="ts">
	import { onMount } from "svelte";
	import { Sse, type SseState } from "../helper/sse";

	type SseResponse<T> = {
		Id: number;
		Path: string;
		Data: T;
	};

	let count: number | undefined;
	let sse: Sse<SseResponse<number>> | undefined = undefined;
	let state: SseState | undefined = undefined;
	let result: string = "";
	let message: string = "";

	onMount(() => {
		start();
	});

	const start = () => {
		sse = new Sse("http://localhost:8080/sse?key=dsadsadsa&topic=notif");
		sse.onData((data) => {
			if (!data) {
				message = "Failed to read data";
				return;
			}
			result = new Date(data.Data).toLocaleString();
		});
		sse.onStateChange((s) => {
			state = s;
		});
		sse.onError((e) => {
			console.log("error", e);
		});
	};

	const stop = () => {
		sse?.close();
		console.log(sse);
		sse = undefined;
		count = undefined;
	};

	const inc = async () => {
		fetch("http://localhost:8080/inc", {
			method: "post",
		});
	};
</script>

<div>
	{#if state}
		<div>
			{state}
		</div>
	{/if}

	<br />
	<div>{result}</div>
	<div>{message}</div>
	<br />
	{#if !sse}
		<button on:click={start}> Start </button>
	{:else}
		<button on:click={stop}> Stop </button>
	{/if}
</div>
