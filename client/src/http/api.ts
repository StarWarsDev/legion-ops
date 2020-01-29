export async function get(uri: string) {
  const res = await fetch(uri);
  const body = await res.json();
  return { status: { code: res.status, text: res.statusText }, body };
}
