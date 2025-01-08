'use client';
import axios from 'axios';
import React, { use, useEffect, useState } from 'react';

interface User {
  ID: string;
  FirstName: string;
  LastName: string;
  Email: string;
}

function Page({ params }: { params: Promise<{ id: string }> }) {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
  const { id } = use(params);
  const [user, setUser] = useState<User>();

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(`${apiUrl}/api/users/${id}`);
        setUser(response.data);
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };
    fetchData();
  });

  return (
    <div>
      <h1>Hello, {user?.FirstName}</h1>
    </div>
  );
}

export default Page;
