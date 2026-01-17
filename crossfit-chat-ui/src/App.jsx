import { useEffect, useMemo, useState } from "react";

const WEEKDAYS = [
  { label: "Segunda", value: "mon" },
  { label: "Terça", value: "tue" },
  { label: "Quarta", value: "wed" },
  { label: "Quinta", value: "thu" },
  { label: "Sexta", value: "fri" },
  { label: "Sábado", value: "sat" },
  { label: "Domingo", value: "sun" },
];

export default function App() {
  const apiUrl = useMemo(
    () => import.meta.env.VITE_API_URL || "http://localhost:8000",
    []
  );

  const [tab, setTab] = useState("chat"); // "chat" | "admin"

  return (
    <div style={styles.page}>
      <div style={styles.container}>
        <header style={styles.header}>
          <div>
            <div style={styles.title}>CrossFit Box</div>
            <div style={styles.subtitle}>
              API: <span style={styles.mono}>{apiUrl}</span>
            </div>
          </div>

          <div style={{ display: "flex", gap: 8 }}>
            <TabButton active={tab === "chat"} onClick={() => setTab("chat")}>
              💬 Chat
            </TabButton>
            <TabButton active={tab === "admin"} onClick={() => setTab("admin")}>
              ⚙️ Admin
            </TabButton>
          </div>
        </header>

        {tab === "chat" ? <Chat apiUrl={apiUrl} /> : <Admin apiUrl={apiUrl} />}
      </div>
    </div>
  );
}

function TabButton({ active, children, onClick }) {
  return (
    <button
      onClick={onClick}
      style={{
        ...styles.tabBtn,
        background: active ? "rgba(11,92,255,0.18)" : "rgba(255,255,255,0.06)",
        border: active
          ? "1px solid rgba(11,92,255,0.35)"
          : "1px solid rgba(255,255,255,0.10)",
      }}
    >
      {children}
    </button>
  );
}

/* =========================================================
   CHAT
========================================================= */
function Chat({ apiUrl }) {
  const [userId] = useState("1");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  const [chat, setChat] = useState([
    {
      from: "bot",
      text: "👋 Olá! Eu sou o assistente do Box. Pergunte sobre planos ou horários 🙂",
      intent: null,
    },
  ]);

  async function sendMessage() {
    const text = message.trim();
    if (!text || loading) return;

    setChat((prev) => [...prev, { from: "user", text }]);
    setMessage("");
    setLoading(true);

    try {
      const res = await fetch(`${apiUrl}/chat`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ user_id: userId, message: text }),
      });

      if (!res.ok) throw new Error(await res.text());

      const data = await res.json();

      setChat((prev) => [
        ...prev,
        {
          from: "bot",
          text: data.reply || "Sem resposta",
          intent: data.intent || null,
        },
      ]);
    } catch (err) {
      setChat((prev) => [
        ...prev,
        {
          from: "bot",
          text: "❌ Erro ao chamar a API. Confira se o backend está rodando.",
          intent: null,
        },
      ]);
      console.error(err);
    } finally {
      setLoading(false);
    }
  }

  function onKeyDown(e) {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  }

  return (
    <>
      <main style={styles.chat}>
        {chat.map((msg, idx) => (
          <div
            key={idx}
            style={{
              ...styles.bubbleRow,
              justifyContent: msg.from === "user" ? "flex-end" : "flex-start",
            }}
          >
            <div
              style={{
                ...styles.bubble,
                background: msg.from === "user" ? "#0b5cff" : "#111827",
              }}
            >
              <div style={styles.bubbleText}>{msg.text}</div>

              {msg.from === "bot" && msg.intent ? (
                <div style={styles.intent}>
                  intent: <span style={styles.mono}>{msg.intent}</span>
                </div>
              ) : null}
            </div>
          </div>
        ))}

        {loading ? (
          <div style={{ ...styles.bubbleRow, justifyContent: "flex-start" }}>
            <div style={{ ...styles.bubble, background: "#111827" }}>
              <div style={styles.bubbleText}>Digitando...</div>
            </div>
          </div>
        ) : null}
      </main>

      <footer style={styles.footer}>
        <textarea
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          onKeyDown={onKeyDown}
          placeholder="Digite uma mensagem (ex: quais são os planos?)"
          style={styles.input}
          rows={2}
        />

        <button
          onClick={sendMessage}
          disabled={loading || !message.trim()}
          style={{
            ...styles.button,
            opacity: loading || !message.trim() ? 0.6 : 1,
            cursor: loading || !message.trim() ? "not-allowed" : "pointer",
          }}
        >
          Enviar
        </button>
      </footer>

      <div style={styles.hints}>
        Sugestões:
        <span style={styles.hintPill}>oi</span>
        <span style={styles.hintPill}>quais são os planos?</span>
        <span style={styles.hintPill}>quais são os horários?</span>
        <span style={styles.hintPill}>sobre o box</span>
      </div>
    </>
  );
}

/* =========================================================
   ADMIN
========================================================= */
function Admin({ apiUrl }) {
  const [section, setSection] = useState("plans"); // "plans" | "schedules"

  return (
    <div style={{ padding: 16, flex: 1, overflow: "auto" }}>
      <div style={{ display: "flex", gap: 8, marginBottom: 12 }}>
        <TabButton active={section === "plans"} onClick={() => setSection("plans")}>
          💳 Planos
        </TabButton>
        <TabButton
          active={section === "schedules"}
          onClick={() => setSection("schedules")}
        >
          ⏰ Horários
        </TabButton>
      </div>

      {section === "plans" ? (
        <AdminPlans apiUrl={apiUrl} />
      ) : (
        <AdminSchedules apiUrl={apiUrl} />
      )}
    </div>
  );
}

function AdminPlans({ apiUrl }) {
  const [plans, setPlans] = useState([]);
  const [loading, setLoading] = useState(false);

  const [form, setForm] = useState({
    name: "",
    description: "",
    price: "",
    active: true,
  });

  async function loadPlans() {
    setLoading(true);
    try {
      const res = await fetch(`${apiUrl}/admin/plans`);
      if (!res.ok) throw new Error(await res.text());
      setPlans(await res.json());
    } catch (err) {
      alert("Erro ao listar planos. Veja console.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  }

  async function createPlan(e) {
    e.preventDefault();

    if (!form.name.trim()) return alert("Informe o nome do plano.");
    if (!String(form.price).trim()) return alert("Informe o preço.");

    const payload = {
      name: form.name.trim(),
      description: form.description.trim() || null,
      price: Number(form.price),
      active: !!form.active,
    };

    try {
      const res = await fetch(`${apiUrl}/admin/plans`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      if (!res.ok) throw new Error(await res.text());

      setForm({ name: "", description: "", price: "", active: true });
      await loadPlans();
    } catch (err) {
      alert("Erro ao criar plano. Veja console.");
      console.error(err);
    }
  }

  useEffect(() => {
    loadPlans();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div style={styles.adminGrid}>
      <div style={styles.card}>
        <div style={styles.cardTitle}>Cadastrar Plano</div>

        <form onSubmit={createPlan} style={{ display: "grid", gap: 10 }}>
          <label style={styles.label}>
            Nome
            <input
              value={form.name}
              onChange={(e) => setForm((p) => ({ ...p, name: e.target.value }))}
              style={styles.input2}
              placeholder="Ex: Livre"
            />
          </label>

          <label style={styles.label}>
            Descrição
            <input
              value={form.description}
              onChange={(e) =>
                setForm((p) => ({ ...p, description: e.target.value }))
              }
              style={styles.input2}
              placeholder="Ex: Acesso ilimitado"
            />
          </label>

          <label style={styles.label}>
            Preço (R$)
            <input
              value={form.price}
              onChange={(e) => setForm((p) => ({ ...p, price: e.target.value }))}
              style={styles.input2}
              placeholder="Ex: 229"
              inputMode="decimal"
            />
          </label>

          <label style={{ ...styles.label, display: "flex", gap: 10 }}>
            <input
              type="checkbox"
              checked={form.active}
              onChange={(e) =>
                setForm((p) => ({ ...p, active: e.target.checked }))
              }
            />
            Ativo
          </label>

          <button style={styles.primaryBtn} type="submit">
            Salvar Plano
          </button>
        </form>
      </div>

      <div style={styles.card}>
        <div style={styles.cardTitle}>
          Planos cadastrados{" "}
          <button onClick={loadPlans} style={styles.smallBtn} disabled={loading}>
            {loading ? "Atualizando..." : "Atualizar"}
          </button>
        </div>

        <div style={styles.tableWrap}>
          <table style={styles.table}>
            <thead>
              <tr>
                <th style={styles.th}>ID</th>
                <th style={styles.th}>Nome</th>
                <th style={styles.th}>Descrição</th>
                <th style={styles.th}>Preço</th>
                <th style={styles.th}>Ativo</th>
              </tr>
            </thead>
            <tbody>
              {plans?.length ? (
                plans.map((p) => (
                  <tr key={p.id}>
                    <td style={styles.td}>{p.id}</td>
                    <td style={styles.td}>{p.name}</td>
                    <td style={styles.td}>{p.description || "-"}</td>
                    <td style={styles.td}>R$ {Number(p.price).toFixed(2)}</td>
                    <td style={styles.td}>{p.active ? "✅" : "❌"}</td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td style={styles.td} colSpan={5}>
                    Nenhum plano cadastrado
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        <div style={styles.tip}>
          💡 Dica: depois de cadastrar, volte na aba <b>Chat</b> e pergunte “quais
          são os planos?”
        </div>
      </div>
    </div>
  );
}

function AdminSchedules({ apiUrl }) {
  const [schedules, setSchedules] = useState([]);
  const [loading, setLoading] = useState(false);

  const [form, setForm] = useState({
    weekday: "mon",
    start_time: "18:00",
    end_time: "19:00",
    modality: "CrossFit",
    coach: "",
    active: true,
  });

  async function loadSchedules() {
    setLoading(true);
    try {
      const res = await fetch(`${apiUrl}/admin/schedules`);
      if (!res.ok) throw new Error(await res.text());
      setSchedules(await res.json());
    } catch (err) {
      alert("Erro ao listar horários. Veja console.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  }

  async function createSchedule(e) {
    e.preventDefault();

    const payload = {
      weekday: form.weekday,
      start_time: form.start_time,
      end_time: form.end_time || null,
      modality: form.modality?.trim() || null,
      coach: form.coach?.trim() || null,
      active: !!form.active,
    };

    try {
      const res = await fetch(`${apiUrl}/admin/schedules`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      if (!res.ok) throw new Error(await res.text());

      setForm((p) => ({ ...p, coach: "" }));
      await loadSchedules();
    } catch (err) {
      alert("Erro ao criar horário. Veja console.");
      console.error(err);
    }
  }

  useEffect(() => {
    loadSchedules();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  function weekdayLabel(val) {
    return WEEKDAYS.find((w) => w.value === val)?.label || val;
  }

  return (
    <div style={styles.adminGrid}>
      <div style={styles.card}>
        <div style={styles.cardTitle}>Cadastrar Horário</div>

        <form onSubmit={createSchedule} style={{ display: "grid", gap: 10 }}>
          <label style={styles.label}>
            Dia da semana
            <select
              value={form.weekday}
              onChange={(e) =>
                setForm((p) => ({ ...p, weekday: e.target.value }))
              }
              style={styles.input2}
            >
              {WEEKDAYS.map((w) => (
                <option key={w.value} value={w.value}>
                  {w.label}
                </option>
              ))}
            </select>
          </label>

          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 10 }}>
            <label style={styles.label}>
              Início
              <input
                type="time"
                value={form.start_time}
                onChange={(e) =>
                  setForm((p) => ({ ...p, start_time: e.target.value }))
                }
                style={styles.input2}
              />
            </label>

            <label style={styles.label}>
              Fim
              <input
                type="time"
                value={form.end_time}
                onChange={(e) =>
                  setForm((p) => ({ ...p, end_time: e.target.value }))
                }
                style={styles.input2}
              />
            </label>
          </div>

          <label style={styles.label}>
            Modalidade
            <input
              value={form.modality}
              onChange={(e) =>
                setForm((p) => ({ ...p, modality: e.target.value }))
              }
              style={styles.input2}
              placeholder="Ex: CrossFit"
            />
          </label>

          <label style={styles.label}>
            Coach (opcional)
            <input
              value={form.coach}
              onChange={(e) => setForm((p) => ({ ...p, coach: e.target.value }))}
              style={styles.input2}
              placeholder="Ex: João"
            />
          </label>

          <label style={{ ...styles.label, display: "flex", gap: 10 }}>
            <input
              type="checkbox"
              checked={form.active}
              onChange={(e) =>
                setForm((p) => ({ ...p, active: e.target.checked }))
              }
            />
            Ativo
          </label>

          <button style={styles.primaryBtn} type="submit">
            Salvar Horário
          </button>
        </form>
      </div>

      <div style={styles.card}>
        <div style={styles.cardTitle}>
          Horários cadastrados{" "}
          <button
            onClick={loadSchedules}
            style={styles.smallBtn}
            disabled={loading}
          >
            {loading ? "Atualizando..." : "Atualizar"}
          </button>
        </div>

        <div style={styles.tableWrap}>
          <table style={styles.table}>
            <thead>
              <tr>
                <th style={styles.th}>ID</th>
                <th style={styles.th}>Dia</th>
                <th style={styles.th}>Início</th>
                <th style={styles.th}>Fim</th>
                <th style={styles.th}>Modalidade</th>
                <th style={styles.th}>Coach</th>
                <th style={styles.th}>Ativo</th>
              </tr>
            </thead>
            <tbody>
              {schedules?.length ? (
                schedules.map((s) => (
                  <tr key={s.id}>
                    <td style={styles.td}>{s.id}</td>
                    <td style={styles.td}>{weekdayLabel(s.weekday)}</td>
                    <td style={styles.td}>{s.start_time?.slice(0, 5)}</td>
                    <td style={styles.td}>
                      {s.end_time ? s.end_time.slice(0, 5) : "-"}
                    </td>
                    <td style={styles.td}>{s.modality || "-"}</td>
                    <td style={styles.td}>{s.coach || "-"}</td>
                    <td style={styles.td}>{s.active ? "✅" : "❌"}</td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td style={styles.td} colSpan={7}>
                    Nenhum horário cadastrado
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        <div style={styles.tip}>
          💡 Dica: depois de cadastrar, volte na aba <b>Chat</b> e pergunte
          “quais são os horários?”
        </div>
      </div>
    </div>
  );
}

/* =========================================================
   STYLES
========================================================= */
const styles = {
  page: {
    height: "100vh",
    background: "#0b1220",
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    padding: 16,
    color: "#fff",
    fontFamily: "Inter, system-ui, -apple-system, Segoe UI, Roboto, Arial",
  },
  container: {
    width: "min(1100px, 100%)",
    height: "min(760px, 100%)",
    background: "#0f172a",
    border: "1px solid rgba(255,255,255,0.06)",
    borderRadius: 16,
    overflow: "hidden",
    display: "flex",
    flexDirection: "column",
    boxShadow: "0 20px 60px rgba(0,0,0,0.4)",
  },
  header: {
    padding: "14px 16px",
    borderBottom: "1px solid rgba(255,255,255,0.06)",
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
  },
  title: { fontSize: 18, fontWeight: 800 },
  subtitle: { fontSize: 12, opacity: 0.8, marginTop: 4 },
  tabBtn: {
    fontSize: 13,
    padding: "8px 12px",
    borderRadius: 999,
    color: "#fff",
  },

  chat: {
    flex: 1,
    padding: 16,
    overflowY: "auto",
    display: "flex",
    flexDirection: "column",
    gap: 10,
  },
  bubbleRow: { display: "flex" },
  bubble: {
    maxWidth: "72%",
    borderRadius: 14,
    padding: "10px 12px",
    lineHeight: 1.35,
  },
  bubbleText: { whiteSpace: "pre-wrap", fontSize: 14 },
  intent: {
    marginTop: 8,
    fontSize: 12,
    opacity: 0.85,
    borderTop: "1px solid rgba(255,255,255,0.08)",
    paddingTop: 8,
  },
  footer: {
    padding: 16,
    borderTop: "1px solid rgba(255,255,255,0.06)",
    display: "flex",
    gap: 10,
  },
  input: {
    flex: 1,
    resize: "none",
    padding: 12,
    borderRadius: 12,
    border: "1px solid rgba(255,255,255,0.10)",
    background: "#0b1220",
    color: "#fff",
    outline: "none",
    fontSize: 14,
  },
  button: {
    padding: "0 18px",
    borderRadius: 12,
    border: "1px solid rgba(255,255,255,0.10)",
    background: "#0b5cff",
    color: "white",
    fontWeight: 800,
  },
  hints: {
    padding: "10px 16px 14px",
    fontSize: 12,
    opacity: 0.85,
    display: "flex",
    gap: 8,
    flexWrap: "wrap",
    borderTop: "1px solid rgba(255,255,255,0.06)",
  },
  hintPill: {
    marginLeft: 6,
    padding: "4px 10px",
    borderRadius: 999,
    background: "rgba(255,255,255,0.06)",
    border: "1px solid rgba(255,255,255,0.08)",
  },

  adminGrid: {
    display: "grid",
    gridTemplateColumns: "1fr 1.4fr",
    gap: 12,
  },
  card: {
    background: "#0b1220",
    border: "1px solid rgba(255,255,255,0.08)",
    borderRadius: 14,
    padding: 14,
  },
  cardTitle: {
    fontSize: 14,
    fontWeight: 800,
    marginBottom: 10,
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
  },
  label: { fontSize: 12, opacity: 0.9, display: "grid", gap: 6 },
  input2: {
    padding: 10,
    borderRadius: 12,
    border: "1px solid rgba(255,255,255,0.10)",
    background: "#0f172a",
    color: "#fff",
    outline: "none",
    fontSize: 14,
  },
  primaryBtn: {
    padding: "10px 12px",
    borderRadius: 12,
    border: "1px solid rgba(11,92,255,0.35)",
    background: "rgba(11,92,255,0.18)",
    color: "#fff",
    fontWeight: 800,
    cursor: "pointer",
  },
  smallBtn: {
    fontSize: 12,
    padding: "6px 10px",
    borderRadius: 999,
    background: "rgba(255,255,255,0.06)",
    border: "1px solid rgba(255,255,255,0.10)",
    color: "#fff",
    cursor: "pointer",
  },
  tableWrap: {
    overflow: "auto",
    borderRadius: 12,
    border: "1px solid rgba(255,255,255,0.08)",
  },
  table: { width: "100%", borderCollapse: "collapse", fontSize: 13 },
  th: {
    textAlign: "left",
    padding: 10,
    fontSize: 12,
    opacity: 0.8,
    background: "rgba(255,255,255,0.03)",
    borderBottom: "1px solid rgba(255,255,255,0.08)",
    position: "sticky",
    top: 0,
  },
  td: {
    padding: 10,
    borderBottom: "1px solid rgba(255,255,255,0.06)",
    verticalAlign: "top",
  },
  tip: {
    marginTop: 10,
    fontSize: 12,
    opacity: 0.9,
    padding: 10,
    borderRadius: 12,
    background: "rgba(255,255,255,0.04)",
    border: "1px solid rgba(255,255,255,0.06)",
  },
  mono: { fontFamily: "ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace" },
};
