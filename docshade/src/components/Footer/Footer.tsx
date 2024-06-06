import React from 'react';
import styled from 'styled-components';
import { FaTelegramPlane } from 'react-icons/fa';

const FooterContainer = styled.footer`
  display: flex;
  align-items: center;
  justify-content: center; /* Center the content */
  position: relative; /* Position relative for absolute positioning */
  padding: 10px 20px;
  background-color: #f0f0f0;
  box-shadow: 0 -4px 8px rgba(0, 0, 0, 0.1);
  width: 100%;

  @media (max-width: 768px) {
    flex-direction: column;
    justify-content: center;
    padding: 20px;
  }
`;

const ContactContainer = styled.div`
  display: flex;
  align-items: center;
  position: absolute; /* Absolute positioning */
  left: 20px; /* Adjust the left position as needed */

  @media (max-width: 768px) {
    position: static;
    margin-bottom: 10px;
  }
`;

const FooterText = styled.p`
  font-size: 14px;
  color: #888;
  text-align: center; /* Center text for mobile view */

  @media (max-width: 768px) {
    text-align: center;
  }
`;

const TelegramLink = styled.a`
  display: flex;
  align-items: center;
  color: #888;
  text-decoration: none;

  &:hover {
    color: #007bff;
  }

  svg {
    margin-right: 8px;
  }
`;

const Footer: React.FC = () => {
  return (
    <FooterContainer>
      <ContactContainer>
        <TelegramLink href="https://t.me/kostyapin" target="_blank" rel="noopener noreferrer">
          <FaTelegramPlane size={24} />
          <FooterText>Contact us on Telegram</FooterText>
        </TelegramLink>
      </ContactContainer>
      <FooterText>DocShade - сервис для анонимизации документов.</FooterText>
    </FooterContainer>
  );
};

export default Footer;
