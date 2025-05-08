'use client'
import { useEffect, useState } from 'react'
import axios from 'axios'
import { Ticket } from '@/gen/ticket_pb'

type Props = {
  id: number
  onClose: () => void
  onUpdated: (t: Ticket) => void
}

export default function TicketModal({ id, onClose, onUpdated }: Props) {
  const [ticket, setTicket] = useState<Ticket | null>(null)
  const [status, setStatus] = useState('')
  const [notes, setNotes] = useState('')

  useEffect(() => {
    axios.get<Ticket>(`http://localhost:8081/tickets/${id}`)
      .then(res => {
        setTicket(res.data)
        setStatus(res.data.status)
        setNotes(res.data.notes)
      })
  }, [id])

  const submit = (e: React.FormEvent) => {
    e.preventDefault()
    axios.patch(`http://localhost:8081/tickets/${id}`, { id, status, notes })
      .then(res => {
        onUpdated(res.data.updatedTicket)
        onClose()
      })
  }

  if (!ticket) return null

  const statusColorMap: Record<string, string> = {
    open: 'bg-gray-200 text-gray-800',
    pending: 'bg-yellow-200 text-yellow-800',
    done: 'bg-green-200 text-green-800'
  }
  
  const statusColor = statusColorMap[ticket.status] ?? 'bg-gray-200 text-gray-800'
  return (
    <div className="fixed inset-0 bg-black/40 z-50 flex items-center justify-center p-4">
      <div className="w-full max-w-lg bg-white p-6 rounded-lg shadow-lg relative animate-fade-in">
        <button
          className="absolute top-2 right-3 text-gray-400 hover:text-black text-2xl"
          onClick={onClose}
        >
          &times;
        </button>

        <div className="mb-4">
          <h2 className="text-2xl font-bold mb-2">{ticket.customerName}</h2>
          <div className={`inline-block px-3 py-1 text-sm rounded ${statusColor}`}>
            {ticket.status}
          </div>
          <div className="text-sm text-gray-500 mt-1">{ticket.email}</div>
          <div className="text-sm text-gray-500">{ticket.createdAt.split('T')[0]}</div>
        </div>

        <form onSubmit={submit} className="space-y-4 mt-4">
          <div>
            <label className="block mb-1 text-sm font-medium">Status</label>
            <select
              className="w-full border border-gray-300 rounded px-3 py-2"
              value={status}
              onChange={e => setStatus(e.target.value)}
            >
              <option value="open">open</option>
              <option value="pending">pending</option>
              <option value="done">done</option>
            </select>
          </div>
          <div>
            <label className="block mb-1 text-sm font-medium">Notes</label>
            <textarea
              className="w-full border border-gray-300 rounded px-3 py-2"
              rows={4}
              value={notes}
              onChange={e => setNotes(e.target.value)}
            />
          </div>
          <div className="flex justify-between">
            <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">Save</button>
            <button type="button" onClick={onClose} className="text-gray-600 hover:underline">Cancel</button>
          </div>
        </form>
      </div>
    </div>
  )
}
