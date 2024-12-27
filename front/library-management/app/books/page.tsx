'use client'

import { useState, useEffect } from 'react'
import { useApi } from '../../hooks/useApi'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

interface Book {
  id: number
  title: string
  author: string
  category: string
}

export default function Books() {
  const [books, setBooks] = useState<Book[]>([])
  const [newBook, setNewBook] = useState({ title: '', author: '', category: '' })
  const [searchTerm, setSearchTerm] = useState('')
  const { fetchApi, loading, error } = useApi()

  useEffect(() => {
    loadBooks()
  }, [])

  const loadBooks = async () => {
    const data = await fetchApi('/books')
    console.log(data); // Логирование данных для проверки
    if (data) {
      // Приведение данных к нужному формату, если API возвращает другие имена
      const formattedBooks = data.map((book: any) => ({
        id: book.ID,
        title: book.Title,
        author: book.Author,
        category: book.Category,
      }))
      setBooks(formattedBooks)
    }
  }

  const addBook = async () => {
    const data = await fetchApi('/books', {
      method: 'POST',
      body: JSON.stringify(newBook)
    })
    if (data) {
      const addedBook = {
        id: data.ID,
        title: data.Title,
        author: data.Author,
        category: data.Category,
      }
      // Обновление состояния с использованием функции, чтобы избежать проблем с асинхронностью
      setBooks((prevBooks) => [...prevBooks, addedBook])
      setNewBook({ title: '', author: '', category: '' })
    }
  }

  const deleteBook = async (id: number) => {
    await fetchApi(`/books/${id}`, { method: 'DELETE' })
    setBooks((prevBooks) => prevBooks.filter(book => book.id !== id))
  }

  const filteredBooks = books.filter(book =>
      (book.title && book.title.toLowerCase().includes(searchTerm.toLowerCase())) ||
      (book.author && book.author.toLowerCase().includes(searchTerm.toLowerCase())) ||
      (book.category && book.category.toLowerCase().includes(searchTerm.toLowerCase()))
  )

  return (
      <div className="space-y-4">
        <h1 className="text-2xl font-bold">Управление книгами</h1>
        <div className="flex space-x-2">
          <Input
              type="text"
              placeholder="Название"
              value={newBook.title}
              onChange={(e) => setNewBook({...newBook, title: e.target.value})}
          />
          <Input
              type="text"
              placeholder="Автор"
              value={newBook.author}
              onChange={(e) => setNewBook({...newBook, author: e.target.value})}
          />
          <Input
              type="text"
              placeholder="Категория"
              value={newBook.category}
              onChange={(e) => setNewBook({...newBook, category: e.target.value})}
          />
          <Button onClick={addBook}>Добавить книгу</Button>
        </div>
        <Input
            type="text"
            placeholder="Поиск по названию, автору или категории"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
        />
        {loading && <p>Загрузка...</p>}
        {error && <p className="text-red-500">{error}</p>}
        {!loading && !error && (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Название</TableHead>
                  <TableHead>Автор</TableHead>
                  <TableHead>Категория</TableHead>
                  <TableHead>Действия</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredBooks.map((book) => (
                    <TableRow key={book.id}>
                      <TableCell>{book.title}</TableCell>
                      <TableCell>{book.author}</TableCell>
                      <TableCell>{book.category}</TableCell>
                      <TableCell>
                        <Button variant="destructive" onClick={() => deleteBook(book.id)}>Удалить</Button>
                      </TableCell>
                    </TableRow>
                ))}
              </TableBody>
            </Table>
        )}
      </div>
  )
}
