import React from 'react';
import styled from 'styled-components';

const HeaderContainer = styled.header`
  display: flex;
  align-items: center;
  justify-content: flex-start; /* Align items to the start */
  padding: 10px 20px;
  background-color: #f0f0f0;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  width: 100%;
  gap: 10px; /* Add a gap between logo and text */

  @media (max-width: 768px) {
    padding: 10px;
    flex-direction: column; /* Stack items on mobile */
    gap: 5px; /* Reduce gap on mobile */
  }
`;

const Logo = styled.img`
  height: 40px;

  @media (max-width: 768px) {
    height: 30px;
  }
`;

const ServiceText = styled.div`
  font-size: 18px;
  color: #333;

  @media (max-width: 768px) {
    text-align: center; /* Center text on mobile */
    font-size: 16px;
  }
`;

const Header: React.FC = () => {
  return (
    <HeaderContainer>
      <Logo src="/logo_mir.png" alt="Logo" />
      <ServiceText>DocShade - сервис для анонимизации документов</ServiceText>
    </HeaderContainer>
  );
};

export default Header;
