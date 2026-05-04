const API_BASE_URL = "http://localhost:8000/api";

interface ApiResponse<T> {
  data?: T;
  error?: string;
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
      console.error(`[API] Error ${response.status}:`, data.error || data.message);
      return { error: data.error || data.message || "An error occurred" };
    }

    // Backend wraps responses in "data" key, so extract it
    const extractedData = data.data !== undefined ? data.data : data;
    console.log(`[API] Extracted data:`, extractedData);
    return { data: extractedData };
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

  verifyResetOTP: (otp: string) =>
    api.post<{ message: string }>("/password_reset/verify_otp", { otp }),

  resetPassword: (password: string, password_confirmation: string) =>
    api.post<{ message: string }>("/password_reset/set_password", {
      password,
      password_confirmation,
    }),
};

