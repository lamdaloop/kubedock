import { useEffect, useState } from 'react';

const API_BASE = 'http://localhost:8080';

export default function App() {
  const [clusters, setClusters] = useState([]);
  const [form, setForm] = useState({ id: '', name: '', url: '', token: '' });
  const [loadingId, setLoadingId] = useState(null);
  const [log, setLog] = useState('');
  const [historyMap, setHistoryMap] = useState({});

  const fetchClusters = async () => {
    const res = await fetch(`${API_BASE}/clusters`);
    const data = await res.json();
    setClusters(data);
  };

  const fetchHistory = async (id) => {
    const res = await fetch(`${API_BASE}/clusters/${id}/history`);
    const data = await res.json();
    setHistoryMap(prev => ({ ...prev, [id]: data }));
  };

  const addCluster = async () => {
    if (!form.id || !form.name || !form.url || !form.token) {
      alert("All fields are required.");
      return;
    }

    await fetch(`${API_BASE}/clusters`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form),
    });

    setForm({ id: '', name: '', url: '', token: '' });
    fetchClusters();
  };

  const triggerBackup = async (id) => {
    setLoadingId(id);
    setLog('');
    const res = await fetch(`${API_BASE}/clusters/${id}/backup`, {
      method: 'POST',
    });
    const text = await res.text();
    setLog(`📦 Backup for '${id}': ${text}`);
    setLoadingId(null);
    fetchHistory(id);
  };
  

  useEffect(() => {
    fetchClusters();
  }, []);

  return (
    <div className="min-h-screen bg-gray-50 text-gray-800">
      <header className="bg-blue-600 text-white p-6 shadow-md">
        <h1 className="text-3xl font-bold text-center">KubeDock Dashboard</h1>
      </header>

      <main className="max-w-5xl mx-auto p-6">
        <section className="bg-white rounded-xl shadow-md p-6 mb-6">
          <h2 className="text-xl font-semibold mb-4">➕ Add New Cluster</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <input className="input" placeholder="Cluster ID" value={form.id} onChange={(e) => setForm({ ...form, id: e.target.value })} />
            <input className="input" placeholder="Cluster Name" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
            <input className="input" placeholder="API Server URL" value={form.url} onChange={(e) => setForm({ ...form, url: e.target.value })} />
            <textarea className="input col-span-1 md:col-span-2" rows={3} placeholder="Bearer Token" value={form.token} onChange={(e) => setForm({ ...form, token: e.target.value })} />
          </div>
          <button className="btn mt-4" onClick={addCluster}>Add Cluster</button>
        </section>

        <section className="bg-white rounded-xl shadow-md p-6 mb-6">
          <h2 className="text-xl font-semibold mb-4">📋 Registered Clusters</h2>
          <table className="w-full table-auto border-collapse mb-4">
            <thead>
              <tr className="bg-gray-100 text-left">
                <th className="p-3 border-b">ID</th>
                <th className="p-3 border-b">Name</th>
                <th className="p-3 border-b">Actions</th>
              </tr>
            </thead>
            <tbody>
              {clusters.map((c, idx) => (
                <tr key={c.id} className={idx % 2 === 0 ? "bg-gray-50" : "bg-white"}>
                  <td className="p-3 border-b">{c.id}</td>
                  <td className="p-3 border-b">{c.name}</td>
                  <td className="p-3 border-b space-x-2">
                  <button className="btn-sm" onClick={() => triggerBackup(c.id)} disabled={loadingId === c.id}>
  {loadingId === c.id ? 'Running...' : 'Backup'}
</button>

                    <button className="btn-outline-sm" onClick={() => fetchHistory(c.id)}>
                      History
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

          {clusters.map(c => historyMap[c.id]?.length > 0 && (
            <div key={c.id} className="mb-6">
              <h3 className="text-lg font-semibold mb-2">🕓 Backup History for {c.name}</h3>
              <table className="w-full text-sm border">
                <thead className="bg-gray-100">
                  <tr>
                    <th className="p-2 border">Time</th>
                    <th className="p-2 border">Status</th>
                    <th className="p-2 border">Path</th>
                  </tr>
                </thead>
                <tbody>
                  {historyMap[c.id].map((h, idx) => (
                    <tr key={h.id} className={idx % 2 === 0 ? "bg-white" : "bg-gray-50"}>
                      <td className="p-2 border">{new Date(h.created_at).toLocaleString()}</td>
                      <td className={`p-2 border font-bold ${h.status === 'success' ? 'text-green-600' : 'text-red-600'}`}>
                        {h.status.toUpperCase()}
                      </td>
                      <td className="p-2 border">{h.path || '-'}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ))}
        </section>

        {log && <div className="p-4 bg-green-100 text-green-800 rounded shadow">{log}</div>}
      </main>
    </div>
  );
}
