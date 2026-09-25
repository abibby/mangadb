const API = "/api";

async function request(path, options) {
  const res = await fetch(API + path, options);
  if (!res.ok) {
    const body = await res.text();
    throw new Error(`${res.status} ${body || res.statusText}`);
  }
  return res.json();
}

export function listSeries(q) {
  const query = q ? `?q=${encodeURIComponent(q)}` : "";
  return request(`/series${query}`).then((r) => r.series);
}

export function getSeries(id, withRelations = []) {
  const query =
    withRelations.length > 0
      ? `?with=${withRelations.map(encodeURIComponent).join("&with=")}`
      : "";
  return request(`/series/${id}${query}`).then((r) => r.series);
}

export function listChapters(id) {
  return request(`/series/${id}/chapters`).then((r) => r.series);
}

export function importSeries(url) {
  return request("/series", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ url }),
  });
}