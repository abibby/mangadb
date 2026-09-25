import { useCallback, useEffect, useState } from "react";
import { getSeries, importSeries, listChapters, listSeries } from "./api.js";
import { Link, Route, Routes, useNavigate, useParams } from "react-router-dom";

function Cover({ series }) {
  if (series.cover_image_url) {
    return <img className="cover" src={series.cover_image_url} alt="" />;
  }
  return <div className="cover cover--placeholder">{series.title}</div>;
}

const SOURCE_URLS = {
  mangaplus: (id) => `https://mangaplus.shueisha.co.jp/titles/${id}`,
  viz: (id) => `https://www.viz.com/shonenjump/chapters/${id}`,
  anilist: (id) => `https://anilist.co/manga/${id}`,
  mangadex: (id) => `https://mangadex.org/title/${id}`,
};

function seriesImportURL(series) {
  const ids = [...(series.ids || [])].sort(
    (a, b) => b.match_quality - a.match_quality
  );
  for (const id of ids) {
    const build = SOURCE_URLS[id.source];
    if (build) return build(id.source_series_id);
  }
  return null;
}

function SeriesCard({ series }) {
  return (
    <Link className="card" to={`/series/${series.id}`}>
      <Cover series={series} />
      <div className="card__title" title={series.title}>
        {series.title}
      </div>
    </Link>
  );
}

function ImportForm() {
  const [url, setUrl] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const navigate = useNavigate();

  async function submit(e) {
    e.preventDefault();
    setBusy(true);
    setError("");
    try {
      const { series_id } = await importSeries(url.trim());
      setUrl("");
      navigate(`/series/${series_id}`);
    } catch (err) {
      setError(err.message);
    } finally {
      setBusy(false);
    }
  }

  return (
    <form className="import" onSubmit={submit}>
      <input
        type="url"
        placeholder="Series URL (e.g. https://mangadex.org/title/...) "
        value={url}
        onChange={(e) => setUrl(e.target.value)}
        required
      />
      <button type="submit" disabled={busy}>
        {busy ? "Importing..." : "Import"}
      </button>
      {error && <span className="error">{error}</span>}
    </form>
  );
}

function SeriesList() {
  const [series, setSeries] = useState(null);
  const [query, setQuery] = useState("");

  const load = useCallback((q) => {
    setSeries(null);
    listSeries(q)
      .then(setSeries)
      .catch((err) => setSeries({ error: err.message }));
  }, []);

  useEffect(() => {
    load("");
  }, [load]);

  function onSearch(e) {
    e.preventDefault();
    load(query.trim());
  }

  return (
    <div>
      <header className="header">
        <div className="header__row">
          <h1>MangaDB</h1>
          <form className="search" onSubmit={onSearch}>
            <input
              type="text"
              placeholder="Search series..."
              value={query}
              onChange={(e) => setQuery(e.target.value)}
            />
            <button type="submit">Search</button>
          </form>
        </div>
        <ImportForm />
      </header>

      {series === null && <p className="notice">Loading...</p>}
      {series?.error && <p className="error">{series.error}</p>}
      {series && !series.error && series.length === 0 && (
        <p className="notice">No series found.</p>
      )}
      {series && !series.error && (
        <div className="grid">
          {series.map((s) => (
            <SeriesCard key={s.id} series={s} />
          ))}
        </div>
      )}
    </div>
  );
}

function ChapterList({ series }) {
  const [chapters, setChapters] = useState(null);
  const [open, setOpen] = useState(null);

  useEffect(() => {
    setChapters(null);
    listChapters(series.id).then(setChapters).catch(() => setChapters([]));
  }, [series.id]);

  const byVolume = {};
  (chapters || []).forEach((ch) => {
    const key = ch.volume ?? 0;
    (byVolume[key] = byVolume[key] || []).push(ch);
  });

  const volumes = [...series.volumes]
    .sort((a, b) => a.number - b.number)
    .map((v) => ({
      ...v,
      chapters: byVolume[v.number] || [],
    }));
  const bareChapters = byVolume[0] || [];

  return (
    <div className="chapters">
      {chapters === null && <p className="notice">Loading chapters...</p>}
      <div className="volumes">
        {volumes.map((v) => {
          const isOpen = open === v.number;
          return (
            <div key={v.number} className="volume">
              <button
                className="volume__card"
                onClick={() => setOpen(isOpen ? null : v.number)}
              >
                {v.cover_image_url ? (
                  <img className="volume__cover" src={v.cover_image_url} alt="" />
                ) : (
                  <div className="volume__cover volume__cover--placeholder">
                    {v.number}
                  </div>
                )}
                <span className="volume__title">Vol. {v.number}</span>
                <span className="volume__count">
                  {v.chapters.length} chapters
                </span>
              </button>
              {isOpen && (
                <ul className="volume__chapters">
                  {v.chapters
                    .slice()
                    .sort((a, b) => a.chapter - b.chapter)
                    .map((ch) => (
                      <li key={ch.chapter}>
                        <span className="chapter-no">
                          Ch. {Number(ch.chapter)}
                        </span>
                        {ch.title || <em>Untitled</em>}
                      </li>
                    ))}
                </ul>
              )}
            </div>
          );
        })}
      </div>
      {bareChapters.length > 0 && (
        <details className="bare">
          <summary>Chapters without a volume ({bareChapters.length})</summary>
          <ul className="volume__chapters">
            {bareChapters
              .slice()
              .sort((a, b) => a.chapter - b.chapter)
              .map((ch) => (
                <li key={ch.chapter}>
                  <span className="chapter-no">Ch. {Number(ch.chapter)}</span>
                  {ch.title || <em>Untitled</em>}
                </li>
              ))}
          </ul>
        </details>
      )}
    </div>
  );
}

function SeriesDetail() {
  const { id } = useParams();
  const [series, setSeries] = useState(null);
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);
  const [refresh, setRefresh] = useState(0);

  useEffect(() => {
    let cancelled = false;
    setSeries(null);
    setError("");

    getSeries(id, ["chapters", "volumes"])
      .then((s) => {
        if (!cancelled) setSeries(s);
      })
      .catch((e) => {
        if (!cancelled) setError(e.message);
      });
    return () => {
      cancelled = true;
    };
  }, [id, refresh]);

  async function reimport() {
    const url = seriesImportURL(series);
    if (!url) {
      setError("No importable source for this series");
      return;
    }
    setBusy(true);
    setError("");
    try {
      await importSeries(url);
      setRefresh((r) => r + 1);
    } catch (err) {
      setError(err.message);
    } finally {
      setBusy(false);
    }
  }

  return (
    <div>
      <div className="toolbar">
        <Link className="back" to="/">
          ← Back
        </Link>
        {series && (
          <button
            className="reimport"
            onClick={reimport}
            disabled={busy || !seriesImportURL(series)}
          >
            {busy ? "Importing..." : "Re-import"}
          </button>
        )}
      </div>
      {error && <p className="error">{error}</p>}
      {!series && !error && <p className="notice">Loading...</p>}
      {series && (
        <>
          <div className="detail">
            <Cover series={series} />
            <div className="detail__info">
              <h1>{series.title}</h1>
              {series.aliases?.length > 0 && (
                <p className="aliases">{series.aliases.join(" / ")}</p>
              )}
              {series.description && (
                <p className="description">{series.description}</p>
              )}
              {series.tags?.length > 0 && (
                <div className="tags">
                  {series.tags.map((t) => (
                    <span key={t} className="tag">
                      {t}
                    </span>
                  ))}
                </div>
              )}
              {series.ids?.length > 0 && (
                <div className="ids">
                  {series.ids.map((id) => (
                    <span key={id.source} className="source">
                      <b>{id.source}</b> {id.source_series_id}
                    </span>
                  ))}
                </div>
              )}
            </div>
          </div>
          <ChapterList key={refresh} series={series} />
        </>
      )}
    </div>
  );
}

export default function App() {
  return (
    <main>
      <Routes>
        <Route path="/" element={<SeriesList />} />
        <Route path="/series/:id" element={<SeriesDetail />} />
      </Routes>
    </main>
  );
}