export interface ServerError {
  message: string;
  errors?: ServerErrorDetail[];
}

export interface ServerErrorDetail {
  field: string;
  message: string;
}
