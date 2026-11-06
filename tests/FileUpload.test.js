import React from 'react';
import { render } from '@testing-library/react';
import FileUpload from '../src/components/FileUpload';

test('renders options by default', () => {
  const { getByText } = render(<FileUpload />);
  expect(getByText(/Language/i)).toBeInTheDocument();
});

test('hides options when showOptions is false', () => {
  const { queryByText } = render(<FileUpload showOptions={false} />);
  expect(queryByText(/Language/i)).not.toBeInTheDocument();
});