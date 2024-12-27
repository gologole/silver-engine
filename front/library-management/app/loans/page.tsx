'use client'

import { useState, useEffect } from 'react'
import { useApi } from '../../hooks/useApi'
import { Button } from "@/components/ui/button"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

interface Loan {
  id: number
  userId: number
  bookId: number
  borrowDate: string
  returnDate: string | null
}

interface User {
  id: number
  name: string
}

interface Book {
  id: number
  title: string
}

export default function Loans() {
  const [loans, setLoans] = useState<Loan[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [books, setBooks] = useState<Book[]>([])
  const [newLoan, setNewLoan] = useState({ userId: '', bookId: '' })
  const { fetchApi, loading, error } = useApi()

  useEffect(() => {
    console.log("Fetching loans, users, and books...");
    loadLoans()
    loadUsers()
    loadBooks()
  }, [])

  const loadLoans = async () => {
    console.log("Loading loans...");
    const data = await fetchApi('/loans')
    console.log("Loans data received:", data);
    if (data) {
      const formattedLoans = data.map((loan: any) => ({
        id: loan.id,                // Приводим к правильному имени поля
        userId: loan.userID,        // Правильное имя поля с сервера
        bookId: loan.bookID,        // Правильное имя поля с сервера
        borrowDate: loan.borrowDate,
        returnDate: loan.returnDate || null,
      }))
      console.log("Formatted loans:", formattedLoans);
      setLoans(formattedLoans)
    }
  }

  const loadUsers = async () => {
    console.log("Loading users...");
    const data = await fetchApi('/users')
    console.log("Users data received:", data);
    if (data) {
      const formattedUsers = data.map((user: any) => ({
        id: user.ID,               // Приводим к правильному имени поля
        name: user.Name,
      }))
      console.log("Formatted users:", formattedUsers);
      setUsers(formattedUsers)
    }
  }

  const loadBooks = async () => {
    console.log("Loading books...");
    const data = await fetchApi('/books')
    console.log("Books data received:", data);
    if (data) {
      const formattedBooks = data.map((book: any) => ({
        id: book.ID,               // Приводим к правильному имени поля
        title: book.Title,
      }))
      console.log("Formatted books:", formattedBooks);
      setBooks(formattedBooks)
    }
  }

  const addLoan = async () => {
    console.log("Adding loan with data:", newLoan);
    const data = await fetchApi('/loans', {
      method: 'POST',
      body: JSON.stringify(newLoan),
    })
    console.log("New loan added:", data);
    if (data) {
      const newLoanData = {
        id: data.ID,
        userId: data.UserID,
        bookId: data.BookID,
        borrowDate: data.BorrowDate,
        returnDate: data.ReturnDate || null,
      }
      setLoans((prevLoans) => [...prevLoans, newLoanData])
      setNewLoan({ userId: '', bookId: '' })
    }
  }

  const returnBook = async (id: number) => {
    console.log(`Returning book with loan ID: ${id}`);
    const data = await fetchApi(`/loans/${id}/return`, { method: 'POST' })
    console.log("Return data received:", data);
    if (data) {
      setLoans((prevLoans) =>
          prevLoans.map((loan) =>
              loan.id === id ? { ...loan, returnDate: data.returnDate } : loan
          )
      )
    }
  }

  return (
      <div className="space-y-4">
        <h1 className="text-2xl font-bold">Выдача и возврат книг</h1>
        <div className="flex space-x-2">
          <Select onValueChange={(value) => {
            console.log("User selected:", value);
            setNewLoan({ ...newLoan, userId: value });
          }}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Выберите пользователя" />
            </SelectTrigger>
            <SelectContent>
              {users.map((user) => (
                  <SelectItem key={user.id} value={user.id.toString()}>
                    {user.name}
                  </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Select onValueChange={(value) => {
            console.log("Book selected:", value);
            setNewLoan({ ...newLoan, bookId: value });
          }}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Выберите книгу" />
            </SelectTrigger>
            <SelectContent>
              {books.map((book) => (
                  <SelectItem key={book.id} value={book.id.toString()}>
                    {book.title}
                  </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Button onClick={addLoan}>Выдать книгу</Button>
        </div>
        {loading && <p>Загрузка...</p>}
        {error && <p className="text-red-500">{error}</p>}
        {!loading && !error && (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID пользователя</TableHead>
                  <TableHead>ID книги</TableHead>
                  <TableHead>Дата выдачи</TableHead>
                  <TableHead>Дата возврата</TableHead>
                  <TableHead>Действия</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {loans.map((loan) => (
                    <TableRow key={`${loan.userId}-${loan.bookId}-${loan.id}`}>
                      <TableCell>{loan.userId}</TableCell>
                      <TableCell>{loan.bookId}</TableCell>
                      <TableCell>{loan.borrowDate}</TableCell>
                      <TableCell>{loan.returnDate || 'Не возвращена'}</TableCell>
                      <TableCell>
                        {!loan.returnDate && (
                            <Button onClick={() => returnBook(loan.id)}>Вернуть книгу</Button>
                        )}
                      </TableCell>
                    </TableRow>
                ))}
              </TableBody>
            </Table>
        )}
      </div>
  )
}
