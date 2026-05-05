const API_BASE_URL = "http://localhost:8000/api";

interface ApiResponse<T> {
  data?: T;
  error?: string;
  success?: string;
}

async function request<T>(
  endpoint: string,
  options: RequestInit = {},
): Promise<ApiResponse<T>> {
  const url = `${API_BASE_URL}${endpoint}`;
  console.log(`[API] ${options.method || "GET"} ${url}`);

  // Get cookie directly
  const cookies = document.cookie;
  console.log(`[API] Current cookies:`, cookies);

  const config: RequestInit = {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    credentials: "include",
  };

  // Also try with explicit cookie header for debugging
  if (cookies.includes("session_token")) {
    console.log(`[API] Session cookie found, should be sent with credentials: include`);
  }

  console.log(`[API] Request config:`, { 
    method: config.method, 
    credentials: config.credentials,
    url: url 
  });

  try {
    const response = await fetch(url, config);
    
    // Log all response headers for debugging
    console.log(`[API] Response status:`, response.status);
    console.log(`[API] Response headers:`, Object.fromEntries(response.headers.entries()));

    const data = await response.json();

    console.log(`[API] Response ${response.status}:`, data);

    if (!response.ok) {
      const errorMsg = data.data?.error || data.error || data.message || "An error occurred";
      console.error(`[API] Error ${response.status}:`, errorMsg);
      return { error: errorMsg };
    }

    // Backend wraps responses in "data" key, so extract it
    const extractedData = data.data !== undefined ? data.data : data;
    console.log(`[API] Extracted data:`, extractedData);
    
    // Check for success message in response
    const successMessage = extractedData?.message;
    return { data: extractedData, success: successMessage };
  } catch (err) {
    console.error("[API] Network error:", err);
    return { error: "Network error" };
  }
}

export const api = {
  get: <T>(endpoint: string) => request<T>(endpoint, { method: "GET" }),
  post: <T>(endpoint: string, body?: unknown) =>
    request<T>(endpoint, {
      method: "POST",
      body: body ? JSON.stringify(body) : undefined,
    }),
  put: <T>(endpoint: string, body?: unknown) =>
    request<T>(endpoint, {
      method: "PUT",
      body: body ? JSON.stringify(body) : undefined,
    }),
  delete: <T>(endpoint: string) => request<T>(endpoint, { method: "DELETE" }),
};

export const authApi = {
  register: (data: {
    name: string;
    email: string;
    password: string;
    password_confirmation: string;
  }) => api.post<{ message: string }>("/auth/register", data),

  login: (data: { email: string; password: string; remember_me?: boolean }) =>
    api.post<{ message: string }>("/auth/login", data),

  logout: () => api.post<{ message: string }>("/auth/logout"),

  getProfile: () =>
    api.get<{ user: { id: string; name: string; email: string } }>("/profile"),

  checkVerified: () =>
    api.get<{ verified: boolean }>("/email_verification/verified"),

  sendOTP: () => api.post<{ message: string }>("/email_verification/send_otp"),

  verifyOTP: (otp: string) =>
    api.post<{ message: string }>("/email_verification/verify_otp", { otp }),

  forgotPassword: (email: string) =>
    api.post<{ message: string }>("/password_reset/send_otp", { email }),

  verifyResetOTP: (email: string, otp: string) =>
    api.post<{ message: string }>("/password_reset/verify_otp", { email, otp }),

  resetPassword: (email: string, password: string, password_confirmation: string) =>
    api.post<{ message: string }>("/password_reset/set_password", {
      email,
      password,
      password_confirmation,
    }),
};

export interface FriendListItem {
  friend_id: string;
  friend_name: string;
  friend_email: string;
  last_message: string;
  last_message_at: string;
  unread_count: number;
}

export interface Friend {
  id: string;
  friend_id: string;
  friend_name: string;
  friend_email: string;
}

export interface MessageListItem {
  message_id: string;
  from_user_id: string;
  to_user_id: string;
  content: string;
  is_me: boolean;
  created_at: string;
}

export interface NotificationListItem {
  id: string;
  type: string;
  content: string;
  is_read: boolean;
  from_user: string;
  created_at: string;
}

export const friendsApi = {
  getFriends: () => api.get<{ friends: FriendListItem[] }>("/friends/list"),

  sendFriendRequest: (to_user_id: string) =>
    api.post<{ message: string }>("/friends/request", { to_user_id }),

  acceptFriendRequest: (request_id: string) =>
    api.post<{ message: string }>("/friends/accept", { request_id }),

  rejectFriendRequest: (request_id: string) =>
    api.post<{ message: string }>("/friends/reject", { request_id }),

  searchFriends: (q: string) =>
    api.get<{ friends: FriendListItem[] }>(`/friends/search?q=${encodeURIComponent(q)}`),

  removeFriend: (friend_id: string) =>
    api.delete<{ message: string }>(`/friends/?friend_id=${friend_id}`),
};

export const searchApi = {
  searchByEmail: (q: string) =>
    api.get<{ users: Friend[] }>(`/search/email?q=${encodeURIComponent(q)}`),

  searchGlobal: (q: string) =>
    api.get<{ messages: MessageListItem[] }>(`/search/global?q=${encodeURIComponent(q)}`),

  searchLocal: (q: string, friend_id: string) =>
    api.get<{ messages: MessageListItem[] }>(
      `/search/local?q=${encodeURIComponent(q)}&friend_id=${friend_id}`
    ),
};

export const chatApi = {
  getMessages: (friend_id: string, limit?: number, offset?: number) => {
    let url = `/chat/messages?friend_id=${friend_id}`;
    if (limit) url += `&limit=${limit}`;
    if (offset) url += `&offset=${offset}`;
    return api.get<{ messages: MessageListItem[] }>(url);
  },

  sendMessage: (to_user_id: string, content: string) =>
    api.post<{ message: MessageListItem }>("/chat/messages", {
      to_user_id,
      content,
    }),
};

export const notificationApi = {
  getNotifications: (limit?: number, offset?: number) => {
    let url = "/notification/list";
    if (limit) url += `?limit=${limit}`;
    if (offset) url += `?offset=${offset}`;
    return api.get<{ notifications: NotificationListItem[] }>(url);
  },

  getUnreadCount: () =>
    api.get<{ unread_count: number }>("/notification/unread-count"),

  markAsRead: (id: string) =>
    api.put<{ message: string }>(`/notification/read?id=${id}`),

  markAllAsRead: () =>
    api.put<{ message: string }>("/notification/read-all"),
};

const WS_URL = "ws://localhost:8000/ws";

export function createWebSocket(userId: string): WebSocket {
  return new WebSocket(`${WS_URL}?user_id=${userId}`);
}

