export async function connect(notificationFunction) {
  const eventSource = new EventSource("/api/events")

  eventSource.onmessage = (e) => {
    notificationFunction(e.data)
  }

  eventSource.onerror = () => {
    console.warn("sse fucked up. it'll probably fix itself.")
  }
}