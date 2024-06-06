import React from 'react';
import StyledButton from './Button.styles';

const Button: React.FC<React.ButtonHTMLAttributes<HTMLButtonElement>> = (props) => (
  <StyledButton {...props} />
);

export default Button;
