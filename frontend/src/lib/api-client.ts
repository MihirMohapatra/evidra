const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

type FetchOptions = RequestInit & {
  params?: Record<string, string>;
};

class ApiError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

async function request<T>(path: string, options: FetchOptions = {}): Promise<T> {
  const { params, ...fetchOpts } = options;

  let url = `${API_URL}${path}`;
  if (params) {
    const qs = new URLSearchParams(params).toString();
    url += `?${qs}`;
  }

  const token = typeof window !== "undefined" ? localStorage.getItem("token") : null;

  const headers: Record<string, string> = {
    ...((fetchOpts.headers as Record<string, string>) || {}),
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  if (!(fetchOpts.body instanceof FormData)) {
    headers["Content-Type"] = "application/json";
  }

  const res = await fetch(url, { ...fetchOpts, headers });

  if (!res.ok) {
    const body = await res.text();
    throw new ApiError(body || res.statusText, res.status);
  }

  if (res.status === 204) {
    return undefined as T;
  }

  return res.json();
}

export const api = {
  // Auth
  login: (email: string, password: string) =>
    request<{ token: string; refresh_token: string; user: any }>("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    }),

  logout: (sessionId: string) =>
    request<void>("/api/v1/auth/logout", {
      method: "POST",
      params: { session_id: sessionId },
    }),

  refresh: (refreshToken: string) =>
    request<{ token: string; refresh_token: string; user: any }>("/api/v1/auth/refresh", {
      method: "POST",
      body: JSON.stringify({ refresh_token: refreshToken }),
    }),

  // Organizations
  listOrganizations: () =>
    request<any[]>("/api/v1/organizations"),

  createOrganization: (name: string, slug: string) =>
    request<any>("/api/v1/organizations", {
      method: "POST",
      body: JSON.stringify({ name, slug }),
    }),

  // Users
  listUsers: (orgId: string) =>
    request<any[]>("/api/v1/users", { params: { organization_id: orgId } }),

  createUser: (data: { organization_id: string; email: string; password: string; role: string }) =>
    request<any>("/api/v1/users", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  // API Keys
  createAPIKey: (orgId: string, name: string) =>
    request<{ key: any; raw_key: string }>("/api/v1/api-keys", {
      method: "POST",
      params: { organization_id: orgId },
      body: JSON.stringify({ name }),
    }),

  revokeAPIKey: (id: string) =>
    request<void>(`/api/v1/api-keys/${id}`, { method: "DELETE" }),

  // Evidence
  listEvidence: (params?: Record<string, string>) =>
    request<{ items: any[]; total: number }>("/api/v1/evidence", { params }),

  getEvidence: (id: string) =>
    request<any>(`/api/v1/evidence/${id}`),

  createEvidence: (data: any) =>
    request<any>("/api/v1/evidence", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  submitEvidence: (id: string) =>
    request<any>(`/api/v1/evidence/${id}/submit`, { method: "POST" }),

  approveEvidence: (id: string, reviewerId: string, comment: string) =>
    request<any>(`/api/v1/evidence/${id}/approve`, {
      method: "POST",
      body: JSON.stringify({ reviewer_id: reviewerId, comment }),
    }),

  rejectEvidence: (id: string, reviewerId: string, comment: string) =>
    request<any>(`/api/v1/evidence/${id}/reject`, {
      method: "POST",
      body: JSON.stringify({ reviewer_id: reviewerId, comment }),
    }),

  // Questionnaires
  listQuestionnaires: () =>
    request<any[]>("/api/v1/questionnaires"),

  // Audit
  listAuditEvents: (params?: Record<string, string>) =>
    request<{ events: any[]; total: number }>("/api/v1/audit/events", { params }),

  // OIDC
  getOIDCProviders: () =>
    request<any[]>("/api/v1/auth/oidc/providers"),
};

export { ApiError };
