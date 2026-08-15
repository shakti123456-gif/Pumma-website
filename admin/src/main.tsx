import { useEffect, useState } from 'react'
import { createRoot } from 'react-dom/client'
import './styles.css'

type Product = { id: number; name: string; description: string; category: string; price: number; stock: number; image: string; active: boolean }
type Dashboard = { products: number; orders: number; revenue: number; lowStock: number }
const api = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'
const categories = ['Shoes', 'Clothing']
const blank: Product = { id: 0, name: '', description: '', category: 'Shoes', price: 0, stock: 0, image: '', active: true }

function App() {
  const [products, setProducts] = useState<Product[]>([])
  const [metrics, setMetrics] = useState<Dashboard>({ products: 0, orders: 0, revenue: 0, lowStock: 0 })
  const [editing, setEditing] = useState<Product | null>(null)
  const [notice, setNotice] = useState('')

  const load = async () => {
    try {
      const [productsResponse, dashboardResponse] = await Promise.all([fetch(`${api}/api/products`), fetch(`${api}/api/admin/dashboard`)])
      setProducts(await productsResponse.json())
      setMetrics(await dashboardResponse.json())
    } catch { setNotice('Cannot reach the API. Start the Gin server at localhost:8080.') }
  }
  useEffect(() => { void load() }, [])
  const save = async (event: React.FormEvent) => {
    event.preventDefault(); if (!editing) return
    const isNew = editing.id === 0
    const response = await fetch(`${api}/api/admin/products${isNew ? '' : `/${editing.id}`}`, { method: isNew ? 'POST' : 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(editing) })
    if (!response.ok) { setNotice('Product could not be saved. Check the form values.'); return }
    setEditing(null); setNotice(isNew ? 'Product created.' : 'Product updated.'); void load()
  }
  const remove = async (id: number) => {
    if (!confirm('Delete this product?')) return
    await fetch(`${api}/api/admin/products/${id}`, { method: 'DELETE' }); setNotice('Product deleted.'); void load()
  }
  const input = (key: keyof Product, label: string, type = 'text') => <label>{label}<input type={type} value={String(editing?.[key] ?? '')} onChange={e => setEditing(p => p ? ({ ...p, [key]: type === 'number' ? Number(e.target.value) : e.target.value }) : p)} /></label>

  return <main>
    <aside><div className="brand">AERO<span>ADMIN</span></div><nav><a className="selected">Overview</a><a href="#products">Products</a><a href="#orders">Orders</a></nav><small>Store management</small></aside>
    <section className="content">
      <header><div><p className="eyebrow">DASHBOARD</p><h1>Good morning, Shakti.</h1><p className="muted">Here is what is happening in your store.</p></div><button onClick={() => setEditing(blank)}>+ Add product</button></header>
      {notice && <p className="notice">{notice}</p>}
      <div className="metrics"><Metric label="Products" value={metrics.products} /><Metric label="Orders" value={metrics.orders} /><Metric label="Revenue" value={`$${metrics.revenue.toFixed(2)}`} /><Metric label="Low stock" value={metrics.lowStock} warning /></div>
      <section className="panel" id="products"><div className="panel-title"><div><p className="eyebrow">CATALOGUE</p><h2>Products</h2></div><button className="quiet" onClick={() => setEditing(blank)}>Add product</button></div>
        <div className="table-wrap"><table><thead><tr><th>Product</th><th>Category</th><th>Price</th><th>Stock</th><th>Status</th><th></th></tr></thead><tbody>{products.map(p => <tr key={p.id}><td><strong>{p.name}</strong><small>{p.description}</small></td><td>{p.category}</td><td>${p.price.toFixed(2)}</td><td className={p.stock < 10 ? 'low' : ''}>{p.stock}</td><td><span className={p.active ? 'badge active' : 'badge'}>{p.active ? 'Active' : 'Hidden'}</span></td><td><button className="text" onClick={() => setEditing(p)}>Edit</button><button className="text danger" onClick={() => void remove(p.id)}>Delete</button></td></tr>)}</tbody></table></div>
      </section>
    </section>
    {editing && <div className="modal-bg"><form className="modal" onSubmit={save}><div className="panel-title"><h2>{editing.id ? 'Edit product' : 'New product'}</h2><button type="button" className="text" onClick={() => setEditing(null)}>Close</button></div><div className="form-grid">{input('name', 'Product name')}<label>Category<select value={editing.category} onChange={e => setEditing({ ...editing, category: e.target.value })}>{!categories.includes(editing.category) && <option value={editing.category}>{editing.category}</option>}{categories.map(category => <option key={category} value={category}>{category}</option>)}</select></label>{input('price', 'Price', 'number')}{input('stock', 'Stock', 'number')}<label className="wide">Description<textarea value={editing.description} onChange={e => setEditing({ ...editing, description: e.target.value })} /></label><label className="toggle"><input type="checkbox" checked={editing.active} onChange={e => setEditing({ ...editing, active: e.target.checked })} /> Available in shop</label></div><div className="actions"><button type="button" className="quiet" onClick={() => setEditing(null)}>Cancel</button><button type="submit">Save product</button></div></form></div>}
  </main>
}
function Metric({ label, value, warning = false }: {label: string; value: string | number; warning?: boolean}) { return <article className="metric"><p>{label}</p><strong className={warning ? 'low' : ''}>{value}</strong><small>{warning ? 'Items need attention' : 'Current total'}</small></article> }
createRoot(document.getElementById('root')!).render(<App />)
