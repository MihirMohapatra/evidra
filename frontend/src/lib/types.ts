export interface Organization {
  id: string;
  name: string;
  slug: string;
  created_at: string;
  updated_at: string;
}

export interface User {
  id: string;
  organization_id: string;
  email: string;
  role: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Session {
  token: string;
  refresh_token: string;
  user: User;
  expires_at: string;
}

export interface APIKey {
  id: string;
  organization_id: string;
  name: string;
  prefix: string;
  is_active: boolean;
  created_at: string;
}

export interface EvidenceItem {
  id: string;
  tenant_id: string;
  title: string;
  content: string;
  category: string;
  source_url: string;
  tags: string[];
  status: string;
  owner_id: string;
  version: number;
  expires_at: string;
  created_at: string;
  updated_at: string;
}

export interface Approval {
  id: string;
  evidence_id: string;
  reviewer_id: string;
  status: string;
  comment: string;
  created_at: string;
}

export interface AuditEvent {
  id: string;
  tenant_id: string;
  actor_id: string;
  action: string;
  target_id: string;
  timestamp: string;
}

export interface Questionnaire {
  id: string;
  tenant_id: string;
  title: string;
  file_url: string;
  file_type: string;
  status: string;
  question_count: number;
  created_at: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  limit: number;
  offset: number;
}
