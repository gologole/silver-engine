'use client'

import { useState, useEffect } from 'react'
import { useApi } from '../../hooks/useApi'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

interface User {
  id: number
  name: string
  email: string
  borrowedBooks: string[]
}

export default function Users() {
  const [users, setUsers] = useState<User[]>([])
  const [newUser, setNewUser] = useState({ name: '', email: '' })
  const [searchTerm, setSearchTerm] = useState('')
  const { fetchApi, loading, error } = useApi()

  useEffect(() => {
    loadUsers()
  }, [])

  const loadUsers = async () => {
    const data = await fetchApi('/users')
    if (data) {
      // Приведение данных с сервера к нужному формату
      const formattedUsers = data.map((user: any) => ({
        id: user.ID,
        name: user.Name,
        email: user.Email,
        borrowedBooks: user.BorrowedBooks || [],
      }))
      setUsers(formattedUsers)
    }
  }

  const addUser = async () => {
    const data = await fetchApi('/users', {
      method: 'POST',
      body: JSON.stringify(newUser)
    })
    if (data) {
      const addedUser = {
        id: data.ID,
        name: data.Name,
        email: data.Email,
        borrowedBooks: data.BorrowedBooks || [],
      }
      // Обновление состояния с использованием функции для корректной работы с асинхронностью
      setUsers((prevUsers) => [...prevUsers, addedUser])
      setNewUser({ name: '', email: '' })
    }
  }

  const deleteUser = async (id: number) => {
    const success = await fetchApi(`/users/${id}`, { method: 'DELETE' })
    if (success) {
      setUsers((prevUsers) => prevUsers.filter(user => user.id !== id))
    }
  }

  const filteredUsers = users.filter(user =>
      (user.name && user.name.toLowerCase().includes(searchTerm.toLowerCase())) ||
      (user.email && user.email.toLowerCase().includes(searchTerm.toLowerCase()))
  )

  return (
      <div className="space-y-4">
        <h1 className="text-2xl font-bold">Управление пользователями</h1>
        <div className="flex space-x-2">
          <Input
              type="text"
              placeholder="Имя"
              value={newUser.name}
              onChange={(e) => setNewUser({ ...newUser, name: e.target.value })}
          />
          <Input
              type="email"
              placeholder="Email"
              value={newUser.email}
              onChange={(e) => setNewUser({ ...newUser, email: e.target.value })}
          />
          <Button onClick={addUser}>Добавить пользователя</Button>
        </div>
        <Input
            type="text"
            placeholder="Поиск по имени или email"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
        />
        {loading && <p>Загрузка...</p>}
        {error && <p className="text-red-500">{error}</p>}
        {!loading && !error && (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Имя</TableHead>
                  <TableHead>Email</TableHead>
                  <TableHead>Выданные книги</TableHead>
                  <TableHead>Действия</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredUsers.map((user) => (
                    <TableRow key={user.id}>
                      <TableCell>{user.name}</TableCell>
                      <TableCell>{user.email}</TableCell>
                      <TableCell>{user.borrowedBooks.join(', ')}</TableCell>
                      <TableCell>
                        <Button variant="destructive" onClick={() => deleteUser(user.id)}>Удалить</Button>
                      </TableCell>
                    </TableRow>
                ))}
              </TableBody>
            </Table>
        )}
      </div>
  )
}
