import { useState, useEffect } from 'react'
import Navbar from '../../components/Navbar'
import { getAllUsers, createUser, updateUser } from '../../api/users'

export default function ManageUsers() {
  const [users, setUsers] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [form, setForm] = useState({ name: '', email: '', password: '', role: 'core' })
  const [creating, setCreating] = useState(false)

  const fetchUsers = async () => {
    try {
      const res = await getAllUsers()
      setUsers(res.data.data || [])
    } catch {
      setError('Failed to load users')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchUsers() }, [])

  const handleChange = (e) => setForm({ ...form, [e.target.name]: e.target.value })

  const handleCreate = async (e) => {
    e.preventDefault()
    setError(''); setSuccess('')
    setCreating(true)
    try {
      await createUser(form)
      setSuccess(`User ${form.email} created successfully`)
      setForm({ name: '', email: '', password: '', role: 'core' })
      fetchUsers()
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to create user')
    } finally {
      setCreating(false)
    }
  }

  const toggleActive = async (user) => {
    try {
      await updateUser(user.id, { is_active: !user.is_active })
      fetchUsers()
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to update user')
    }
  }

  return (
    <>
      <Navbar />
      <div className="page-container">
        <h2>Manage Users</h2>

        <div className="form-card">
          <h3>Create New User</h3>
          {error && <div className="alert alert-error">{error}</div>}
          {success && <div className="alert alert-success">{success}</div>}
          <form onSubmit={handleCreate}>
            <div className="form-grid">
              <div className="form-group">
                <label>Full Name</label>
                <input type="text" name="name" value={form.name} onChange={handleChange} required placeholder="Full name" />
              </div>
              <div className="form-group">
                <label>Email</label>
                <input type="email" name="email" value={form.email} onChange={handleChange} required placeholder="user@iith.ac.in" />
              </div>
              <div className="form-group">
                <label>Password</label>
                <input type="password" name="password" value={form.password} onChange={handleChange} required placeholder="Min 6 characters" minLength={6} />
              </div>
              <div className="form-group">
                <label>Role</label>
                <select name="role" value={form.role} onChange={handleChange}>
                  <option value="core">Core</option>
                  <option value="viewer">Viewer</option>
                  <option value="admin">Admin</option>
                </select>
              </div>
            </div>
            <button type="submit" className="btn btn-primary" disabled={creating}>
              {creating ? 'Creating...' : 'Create User'}
            </button>
          </form>
        </div>

        <div className="section-header"><h3>All Users</h3></div>
        {loading ? (
          <div className="loading">Loading users...</div>
        ) : (
          <div className="table-wrapper">
            <table className="data-table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Email</th>
                  <th>Role</th>
                  <th>Status</th>
                  <th>Created</th>
                  <th>Action</th>
                </tr>
              </thead>
              <tbody>
                {users.map(u => (
                  <tr key={u.id}>
                    <td>{u.name}</td>
                    <td>{u.email}</td>
                    <td><span className={`role-badge role-${u.role}`}>{u.role}</span></td>
                    <td>
                      <span className={`status-badge ${u.is_active ? 'confirmed' : 'cancelled'}`}>
                        {u.is_active ? 'Active' : 'Inactive'}
                      </span>
                    </td>
                    <td>{new Date(u.created_at).toLocaleDateString()}</td>
                    <td>
                      <button
                        className={`btn btn-sm ${u.is_active ? 'btn-danger' : 'btn-success'}`}
                        onClick={() => toggleActive(u)}
                      >
                        {u.is_active ? 'Deactivate' : 'Activate'}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </>
  )
}
