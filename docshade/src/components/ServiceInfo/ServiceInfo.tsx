import React from 'react';
import styled from 'styled-components';

const ServiceInfoContainer = styled.div`
  max-width: 600px;
  margin: 40px auto;
  padding: 20px;
  background: #f9f9f9;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
`;

const Title = styled.h2`
  font-size: 24px;
  margin-bottom: 20px;
`;

const Subtitle = styled.p`
  font-size: 16px;
  color: #888;
  margin-bottom: 20px;
`;

const StepsContainer = styled.div`
  display: flex;
  justify-content: space-around;
  flex-wrap: wrap;
  text-align: left;
`;

const Step = styled.div`
  flex: 1;
  margin: 20px 10px;
  min-width: 150px;
`;

const StepNumber = styled.div`
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 10px;
  text-align: center;
`;

const StepDescription = styled.p`
  font-size: 14px;
  color: #333;
  text-align: center;
`;

const ImportantNote = styled.div`
  margin-top: 20px;
  font-size: 14px;
  color: #555;
  text-align: center;
`;

const ServiceInfo: React.FC = () => {
  return (
    <ServiceInfoContainer>
      <Title>Как обезличить текст и персональные данные из файлов PDF</Title>
      <Subtitle>Как это работает</Subtitle>
      <StepsContainer>
        <Step>
          <StepNumber>ШАГ 1</StepNumber>
          <StepDescription>Кликните по области загрузки или просто перетащите туда ваш файл pdf.</StepDescription>
        </Step>
        <Step>
          <StepNumber>ШАГ 2</StepNumber>
          <StepDescription>После обработки файла нажмите кнопку загрузки, если она не началась автоматически.</StepDescription>
        </Step>
        <Step>
          <StepNumber>ШАГ 3</StepNumber>
          <StepDescription>Обработанные документы хранятся в хранилище вашего браузера и доступны до повторного сохранения</StepDescription>
        </Step>
      </StepsContainer>
      <ImportantNote>
        Важно: сервис не хранит документы пользователя и беспокоится о их безопасности и анонимности в сети.
      </ImportantNote>
    </ServiceInfoContainer>
  );
};

export default ServiceInfo;
